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
	GetAllOrders(ctx context.Context) (order []*model.Order, err error)
}

type OrdersServiceImpl struct {
	table db.OrdersTable
	cache *bigcache.BigCache
}

func newOrdersService(ordersTable db.OrdersTable) OrdersService {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         0,
		CleanWindow:        0,
		MaxEntriesInWindow: 1000 * 10 * 60,
		HardMaxCacheSize:   1,
	}

	cache, initErr := bigcache.NewBigCache(config)
	if initErr != nil {
		log.Fatal(initErr)
	}

	ordersService := &OrdersServiceImpl{ordersTable, cache}

	orders, err := ordersService.table.GetAllOrders(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {
		err = ordersService.saveOrderInCache(order)
		if err != nil {
			log.Fatal(err)
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

func (service *OrdersServiceImpl) GetAllOrders(ctx context.Context) (order []*model.Order, err error) {
	return service.table.GetAllOrders(ctx)
}
