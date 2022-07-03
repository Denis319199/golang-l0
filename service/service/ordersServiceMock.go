package service

import (
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"service/db"
	"service/model"
)

const testOrderJson = `
{
"order_uid": "1",
"track_number": "WBILMTESTTRACK",
"entry": "WBIL",
"delivery": {
  "name": "Test Testov",
  "phone": "+9720000000",
  "zip": "2639809",
  "city": "Kiryat Mozkin",
  "address": "Ploshad Mira 15",
  "region": "Kraiot",
  "email": "test@gmail.com"
},
"payment": {
  "transaction": "b563feb7b2b84b6test",
  "request_id": "",
  "currency": "USD",
  "provider": "wbpay",
  "amount": 1817,
  "payment_dt": 1637907727,
  "bank": "alpha",
  "delivery_cost": 1500,
  "goods_total": 317,
  "custom_fee": 0
},
"items": [
  {
	"chrt_id": 9934930,
	"track_number": "WBILMTESTTRACK",
	"price": 453,
	"rid": "ab4219087a764ae0btest",
	"name": "Mascaras",
	"sale": 30,
	"size": "0",
	"total_price": 317,
	"nm_id": 2389212,
	"brand": "Vivienne Sabo",
	"status": 202
  }
],
"locale": "en",
"internal_signature": "",
"customer_id": "test",
"delivery_service": "meest",
"shardkey": "9",
"sm_id": 99,
"date_created": "2021-11-26T06:22:19Z",
"oof_shard": "1"
}`

type mockError struct {
}

func (err mockError) Error() string {
	return "mock error"
}

type OrdersTableMock struct {
	returnError bool
}

func (table OrdersTableMock) InsertOrder(ctx context.Context, order *model.Order) error {
	if table.returnError {
		return mockError{}
	}

	return nil
}

func (table OrdersTableMock) GetOrderById(ctx context.Context, id string) (res *model.Order, err error) {
	if table.returnError {
		return nil, mockError{}
	}

	order := model.Order{}
	err = json.Unmarshal([]byte(testOrderJson), &order)
	if err != nil {
		log.Fatal(err)
	}

	order.OrderUid = &id
	return &order, nil
}

func (table OrdersTableMock) GetAllOrders(ctx context.Context) (res []*model.Order, err error) {
	if table.returnError {
		return nil, mockError{}
	}

	order1 := model.Order{}
	err = json.Unmarshal([]byte(testOrderJson), &order1)
	if err != nil {
		log.Fatal(err)
	}

	orderUid1 := "1"
	orderUid2 := "2"
	orderUid3 := "3"
	orderUid4 := "4"
	orderUid5 := "5"

	order1.OrderUid = &orderUid1
	order2, order3, order4, order5 := order1, order1, order1, order1
	order2.OrderUid = &orderUid2
	order3.OrderUid = &orderUid3
	order4.OrderUid = &orderUid4
	order5.OrderUid = &orderUid5

	orders := []*model.Order{&order1, &order2, &order3, &order4, &order5}
	return orders, nil
}

func MockOrdersTable(returnError bool) db.OrdersTable {
	return OrdersTableMock{returnError}
}

func CreateOrdersServiceWithCache(returnError bool) *OrdersServiceImpl {
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

	return &OrdersServiceImpl{MockOrdersTable(returnError), cache}
}
