package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"service/model"
)

const (
	getAllOrdersQuery = "SELECT * FROM orders;"
	getOrderByIdQuery = "SELECT * FROM orders WHERE order_uid = $1;"
	insertIntoQuery   = "INSERT INTO orders VALUES($1, $2);"
)

type OrdersTable interface {
	GetAllOrders(ctx context.Context) (res []*model.Order, err error)
	GetOrderById(ctx context.Context, id string) (res *model.Order, err error)
	InsertOrder(ctx context.Context, order *model.Order) error
}

type OrdersTableImpl struct {
	database              *sql.DB
	getAllOrdersStatement *sql.Stmt
	getOrderByIdStatement *sql.Stmt
	insertIntoStatement   *sql.Stmt
}

func newOrdersTable(db *sql.DB) OrdersTable {
	table := OrdersTableImpl{database: db}

	var err error
	table.getAllOrdersStatement, err = db.Prepare(getAllOrdersQuery)
	if err != nil {
		log.Fatal(err)
	}

	table.getOrderByIdStatement, err = db.Prepare(getOrderByIdQuery)
	if err != nil {
		log.Fatal(err)
	}

	table.insertIntoStatement, err = db.Prepare(insertIntoQuery)
	if err != nil {
		log.Fatal(err)
	}

	return &table
}

func (table *OrdersTableImpl) GetAllOrders(ctx context.Context) (res []*model.Order, err error) {
	var rows *sql.Rows
	rows, err = table.getAllOrdersStatement.QueryContext(ctx)
	if err != nil {
		return
	}

	for rows.Next() {
		var order *model.Order
		order, err = mapToOrder(rows)

		if err != nil {
			return
		}

		res = append(res, order)
	}

	return
}

func (table *OrdersTableImpl) GetOrderById(ctx context.Context, id string) (res *model.Order, err error) {
	var rows *sql.Rows
	rows, err = table.getOrderByIdStatement.QueryContext(ctx, id)

	if rows.Next() {
		res, err = mapToOrder(rows)
	} else {
		err = NotFoundError{}
	}

	return
}

func (table *OrdersTableImpl) InsertOrder(ctx context.Context, order *model.Order) error {
	marshal, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, err = table.insertIntoStatement.ExecContext(ctx, order.OrderUid, marshal)
	if err != nil {
		return err
	}

	return nil
}

func mapToOrder(rows *sql.Rows) (res *model.Order, err error) {
	var (
		orderId  string
		orderRow []byte
	)

	err = rows.Scan(&orderId, &orderRow)
	if err != nil {
		return
	}

	err = json.Unmarshal(orderRow, &res)
	return
}
