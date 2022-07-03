package config

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"log"
	"service/app"
	"service/model"
)

func newNatsStreaming() stan.Conn {
	viper.SetDefault("nats.cluster-id", "test-cluster")
	viper.SetDefault("nats.client-id", "service")
	viper.SetDefault("nats.url", "nats://127.0.0.1:4222")

	clusterId := viper.GetString("nats.cluster-id")
	clientId := viper.GetString("nats.client-id")
	url := viper.GetString("nats.url")
	if string([]rune(url)[:7]) != "nats://" {
		url = "nats://" + url
	}

	nats, err := stan.Connect(
		clusterId,
		clientId,
		stan.NatsURL(url),
	)
	if err != nil {
		log.Fatalln(err)
	}

	setupNatsHandlers(nats)

	return nats
}

func setupNatsHandlers(nats stan.Conn) {
	_, err := nats.Subscribe("myevent", processMyEvent)
	if err != nil {
		log.Fatalln(err)
	}
}

func processMyEvent(m *stan.Msg) {
	order := &model.Order{}
	err := json.Unmarshal(m.Data, order)
	if err != nil {
		return
	}

	err = app.ServiceInstance.Validator.Struct(order)
	if err != nil {
		return
	}

	_ = app.ServiceInstance.Services.OrdersService.InsertOrder(context.Background(), order)
}
