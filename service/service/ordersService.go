package service

import (
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"service/db"
	"service/model"
)

type OrdersService interface {
	InsertOrder(ctx context.Context, order *model.Order) error
	GetOrderById(ctx context.Context, orderUid string) (order *model.Order, err error)
	GetOrders(page, size int, ctx context.Context) (order []*model.Order, err error)
}

type OrdersServiceImpl struct {
	table db.OrdersTable
	cache *bigcache.BigCache
}

func newOrdersService(ordersTable db.OrdersTable) OrdersService {
	const maxCacheSize = 5

	config := bigcache.Config{
		Shards:             512,
		LifeWindow:         0,
		CleanWindow:        0,
		MaxEntriesInWindow: 100000,
		HardMaxCacheSize:   maxCacheSize,
	}

	cache, initErr := bigcache.NewBigCache(config)
	if initErr != nil {
		log.Fatal(initErr)
	}

	ordersService := &OrdersServiceImpl{ordersTable, cache}

	// Fill cache with orders until cache size exceeds 4MB (max cache size 5MB)
	for page := 0; cache.Capacity() < ((maxCacheSize - 1) << 20); page++ {
		orders, err := ordersService.table.GetOrders(page, 50, context.Background())
		if err != nil {
			log.Fatal(err)
		}

		if orders == nil {
			break
		}

		for _, order := range orders {
			err = ordersService.saveOrderInCache(order)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return ordersService
}

func (service *OrdersServiceImpl) saveOrderInCache(order *model.Order) error {
	marshal, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return service.cache.Set(*order.OrderUid, marshal)
}

func (service *OrdersServiceImpl) searchOrderInCacheById(orderUid string) (order *model.Order, err error) {
	var orderRow []byte

	orderRow, err = service.cache.Get(orderUid)
	if err != nil {
		return
	}

	order = &model.Order{}
	err = json.Unmarshal(orderRow, order)
	if err != nil {
		return
	}

	return
}

func (service *OrdersServiceImpl) InsertOrder(ctx context.Context, order *model.Order) error {
	err := service.table.InsertOrder(ctx, order)
	if err == nil {
		err := service.saveOrderInCache(order)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func (service *OrdersServiceImpl) GetOrderById(ctx context.Context, orderUid string) (order *model.Order, err error) {
	order, err = service.searchOrderInCacheById(orderUid)
	if err == nil {
		return
	}

	order, err = service.table.GetOrderById(ctx, orderUid)
	if err == nil {
		err := service.saveOrderInCache(order)
		if err != nil {
			log.Println(err)
		}
	}

	return order, err
}

func (service *OrdersServiceImpl) GetOrders(page, size int, ctx context.Context) (order []*model.Order, err error) {
	return service.table.GetOrders(page, size, ctx)
}
