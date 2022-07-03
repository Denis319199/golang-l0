package service

import "service/db"

type Services struct {
	OrdersService OrdersService
}

func NewServices(database *db.Database) *Services {
	return &Services{newOrdersService(database.OrdersTable)}
}
