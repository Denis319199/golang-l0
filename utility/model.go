package main

import "time"

type Delivery struct {
	Name    *string `json:"name" validate:"required""`
	Phone   *string `json:"phone" validate:"required"`
	Zip     *string `json:"zip" validate:"required"`
	City    *string `json:"city" validate:"required"`
	Address *string `json:"address" validate:"required"`
	Region  *string `json:"region" validate:"required"`
	Email   *string `json:"email" validate:"required"`
}

type Payment struct {
	Transaction  *string `json:"transaction" validate:"required"`
	RequestId    *string `json:"request_id" validate:"required"`
	Currency     *string `json:"currency" validate:"required"`
	Provider     *string `json:"provider" validate:"required"`
	Amount       *uint   `json:"amount" validate:"required"`
	PaymentDt    *uint64 `json:"payment_dt" validate:"required"`
	Bank         *string `json:"bank"`
	DeliveryCost *uint   `json:"delivery_cost" validate:"required"`
	GoodsTotal   *uint   `json:"goods_total" validate:"required"`
	CustomFee    *uint   `json:"custom_fee" validate:"required,min=0"`
}

type Item struct {
	ChrtId      *uint   `json:"chrt_id" validate:"required"`
	TrackNumber *string `json:"track_number" validate:"required"`
	Price       *uint   `json:"price" validate:"required"`
	Rid         *string `json:"rid" validate:"required"`
	Name        *string `json:"name" validate:"required"`
	Sale        *uint   `json:"sale" validate:"required"`
	Size        *string `json:"size" validate:"required"`
	TotalPrice  *uint   `json:"total_price" validate:"required"`
	NmId        *uint   `json:"nm_id" validate:"required"`
	Brand       *string `json:"brand" validate:"required"`
	Status      *uint   `json:"status" validate:"required"`
}

type Order struct {
	OrderUid          *string    `json:"order_uid" validate:"required"`
	TrackNumber       *string    `json:"track_number" validate:"required"`
	Entry             *string    `json:"entry" validate:"required"`
	Delivery          *Delivery  `json:"delivery" validate:"required"`
	Payment           *Payment   `json:"payment" validate:"required"`
	Items             *[]Item    `json:"items" validate:"required"`
	Locale            *string    `json:"locale" validate:"required"`
	InternalSignature *string    `json:"internal_signature" validate:"required"`
	CustomerId        *string    `json:"customer_id" validate:"required"`
	DeliveryService   *string    `json:"delivery_service" validate:"required"`
	Shardkey          *string    `json:"shardkey" validate:"required"`
	SmId              *uint      `json:"sm_id" validate:"required"`
	DateCreated       *time.Time `json:"date_created" validate:"required"`
	OofShard          *string    `json:"oof_shard" validate:"required"`
}
