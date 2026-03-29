package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/redis/go-redis/v9"
)

var _ domain.CacheRepository = (*Cache)(nil)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (r *Cache) getKey(productID int) string {
	if productID == -1 {
		return "product:*:packs"
	}

	return fmt.Sprintf("product:%d:packs", productID)
}

func (r *Cache) GetPacksByProductIDs(ctx context.Context, productIDs []int) ([]*domain.ProductPacks, error) {
	pipe := r.client.Pipeline()

	cmds := make(map[int]*redis.StringCmd, len(productIDs))
	for _, id := range productIDs {
		cmds[id] = pipe.Get(ctx, r.getKey(id))
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		allNil := true
		for _, cmd := range cmds {
			if cmd.Err() != nil && cmd.Err() != redis.Nil {
				allNil = false
				break
			}
		}
		if !allNil {
			return nil, err
		}
	}

	var result []*domain.ProductPacks
	for _, id := range productIDs {
		cmd := cmds[id]
		val, err := cmd.Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				result = append(result, &domain.ProductPacks{ProductID: id, Packs: []*domain.Pack{}})
				continue
			}

			return nil, err
		}

		var packs []*domain.Pack
		if err := json.Unmarshal([]byte(val), &packs); err != nil {
			return nil, err
		}

		result = append(result, &domain.ProductPacks{ProductID: id, Packs: packs})
	}

	return result, nil
}

func (r *Cache) SetPack(ctx context.Context, pack *domain.Pack) error {
	key := r.getKey(pack.ProductID)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	var packs []*domain.Pack
	if err == nil {
		if err := json.Unmarshal([]byte(val), &packs); err != nil {
			return err
		}
	}

	found := false
	for i, p := range packs {
		if p.ID == pack.ID {
			packs[i] = pack
			found = true
			break
		}
	}

	if !found {
		packs = append(packs, pack)
	}

	data, err := json.Marshal(packs)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *Cache) SetPacks(ctx context.Context, packs []*domain.Pack) error {
	grouped := make(map[int][]*domain.Pack)
	for _, p := range packs {
		grouped[p.ProductID] = append(grouped[p.ProductID], p)
	}

	pipe := r.client.Pipeline()
	for productID, productPacks := range grouped {
		data, err := json.Marshal(productPacks)
		if err != nil {
			return err
		}
		pipe.Set(ctx, r.getKey(productID), data, 0)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Cache) RemovePack(ctx context.Context, id int) error {
	var cursor uint64

	for {
		keys, nextCursor, err := r.client.
			Scan(ctx, cursor, r.getKey(-1), 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			val, err := r.client.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var packs []*domain.Pack
			if err := json.Unmarshal([]byte(val), &packs); err != nil {
				continue
			}

			filtered := make([]*domain.Pack, 0, len(packs))
			for _, p := range packs {
				if p.ID != id {
					filtered = append(filtered, p)
				}
			}

			if len(filtered) != len(packs) {
				data, err := json.Marshal(filtered)
				if err != nil {
					return err
				}

				if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
					return err
				}

				return nil
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *Cache) Invalidate(ctx context.Context) error {
	var cursor uint64

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, r.getKey(-1), 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
