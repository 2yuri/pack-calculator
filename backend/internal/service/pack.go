package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/2yuri/pack-calculator/internal/domain"
)

type PackDeps struct {
	Repo      domain.PackRepository
	CacheRepo domain.CacheRepository
}

var _ domain.PackService = (*Pack)(nil)

type Pack struct {
	repo      domain.PackRepository
	cacheRepo domain.CacheRepository
}

func NewPack(deps PackDeps) *Pack {
	return &Pack{
		repo:      deps.Repo,
		cacheRepo: deps.CacheRepo,
	}
}

func (s *Pack) SyncCache(ctx context.Context) error {
	packs, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	if err := s.cacheRepo.Invalidate(ctx); err != nil {
		return err
	}

	if len(packs) > 0 {
		if err := s.cacheRepo.SetPacks(ctx, packs); err != nil {
			return err
		}
	}

	zap.L().Info("cache synced from database", zap.Int("packs", len(packs)))

	return nil
}

func (s *Pack) GetByProductID(ctx context.Context, productID int) ([]*domain.Pack, error) {
	return s.repo.GetByProductID(ctx, productID)
}

func (s *Pack) Create(ctx context.Context, productID int, size int) (*domain.Pack, error) {
	pack, err := s.repo.Create(ctx, productID, size)
	if err != nil {
		return nil, err
	}

	if err := s.cacheRepo.SetPack(ctx, pack); err != nil {
		zap.L().Error("failed to update cache after create", zap.Error(err))
	}

	return pack, nil
}

func (s *Pack) CreateBatch(ctx context.Context, productID int, sizes []int) ([]*domain.Pack, error) {
	packs, err := s.repo.CreateBatch(ctx, productID, sizes)
	if err != nil {
		return nil, err
	}

	if err := s.cacheRepo.SetPacks(ctx, packs); err != nil {
		zap.L().Error("failed to update cache after batch create", zap.Error(err))
	}

	return packs, nil
}

func (s *Pack) Update(ctx context.Context, id int, size int) (*domain.Pack, error) {
	pack, err := s.repo.Update(ctx, id, size)
	if err != nil {
		return nil, err
	}

	if err := s.cacheRepo.SetPack(ctx, pack); err != nil {
		zap.L().Error("failed to update cache after update", zap.Error(err))
	}

	return pack, nil
}

func (s *Pack) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.cacheRepo.RemovePack(ctx, id); err != nil {
		zap.L().Error("failed to update cache after delete", zap.Error(err))
	}

	return nil
}
