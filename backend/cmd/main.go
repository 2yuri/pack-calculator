package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/2yuri/pack-calculator/internal/api"
	orderhandler "github.com/2yuri/pack-calculator/internal/api/handler/order"
	packhandler "github.com/2yuri/pack-calculator/internal/api/handler/pack"
	producthandler "github.com/2yuri/pack-calculator/internal/api/handler/product"
	swaggerhandler "github.com/2yuri/pack-calculator/internal/api/handler/swagger"
	"github.com/2yuri/pack-calculator/internal/repository"
	"github.com/2yuri/pack-calculator/internal/service"
	"github.com/2yuri/pack-calculator/pkg/cache"
	"github.com/2yuri/pack-calculator/pkg/config"
	"github.com/2yuri/pack-calculator/pkg/db"
	"github.com/2yuri/pack-calculator/pkg/rest"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	handleErr(err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	conn, err := db.New()
	handleErr(err)
	defer conn.Close()

	redisClient, err := cache.New()
	handleErr(err)
	defer redisClient.Close()

	productRepo := repository.NewProduct(conn)
	packRepo := repository.NewPack(conn)
	cacheRepo := repository.NewCache(redisClient)

	productSvc := service.NewProduct(service.ProductDeps{
		Repo: productRepo,
	})
	packSvc := service.NewPack(service.PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	orderSvc := service.NewOrder(service.OrderDeps{
		CacheRepo: cacheRepo,
	})

	handleErr(packSvc.SyncCache(ctx))

	server := api.New(api.Deps{
		Port:   config.Instance().App.Port,
		Prefix: "/api/v1",
		Routes: []rest.Routes{
			producthandler.NewHandler(producthandler.Deps{Service: productSvc}),
			packhandler.NewHandler(packhandler.Deps{Service: packSvc}),
			orderhandler.NewHandler(orderhandler.Deps{Service: orderSvc}),
			swaggerhandler.NewHandler(),
		},
	})

	if err := server.Start(); err != nil {
		zap.L().Fatal("server failed", zap.Error(err))
	}
}
