package service

import (
	"context"
	"encoding/json"
	"errors"
	"service/model"
	"testing"
)

func TestCachingInsertOrder(t *testing.T) {
	service := CreateOrdersServiceWithCache(false)

	orderUid := "1"
	order := model.Order{OrderUid: &orderUid}
	err := service.InsertOrder(context.Background(), &order)
	if err != nil {
		t.Error(err)
	}

	var orderRowGot []byte
	orderRowGot, err = service.cache.Get(orderUid)
	if err != nil {
		t.Error(err)
	}

	var orderGot model.Order
	err = json.Unmarshal(orderRowGot, &orderGot)
	if err != nil {
		t.Error(err)
	}

	if *orderGot.OrderUid != orderUid {
		t.Error("Saved orderUid (", *orderGot.OrderUid,
			") and given orderUid (", orderUid, ") are not equal")
	}
}

func TestCachingInsertOrderWithError(t *testing.T) {
	service := CreateOrdersServiceWithCache(true)

	orderUid := "1"
	order := model.Order{OrderUid: &orderUid}
	err := service.InsertOrder(context.Background(), &order)
	if err != nil && !errors.Is(err, mockError{}) {
		t.Error(err)
	}

	_, err = service.cache.Get(orderUid)
	if err == nil {
		t.Error(err)
	}
}

func TestCachingGetOrderById(t *testing.T) {
	service := CreateOrdersServiceWithCache(false)

	orderUid := "1"
	order, err := service.GetOrderById(context.Background(), orderUid)
	if err != nil {
		t.Error(err)
	}

	var orderRowGot []byte
	orderRowGot, err = service.cache.Get(orderUid)
	if err != nil {
		t.Error(err)
	}

	var orderGot model.Order
	err = json.Unmarshal(orderRowGot, &orderGot)
	if err != nil {
		t.Error(err)
	}

	if *orderGot.OrderUid != *order.OrderUid {
		t.Error("Saved orderUid (", *orderGot.OrderUid,
			") and given orderUid (", *order.OrderUid, ") are not equal")
	}
}
