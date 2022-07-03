package main

import (
	"github.com/nats-io/stan.go"
)

func setupNatsStreaming(url, clusterId string) stan.Conn {
	clientId := "this-service"
	if string([]rune(url)[:7]) != "nats://" {
		url = "nats://" + url
	}

	nats, err := stan.Connect(
		clusterId,
		clientId,
		stan.NatsURL(url),
	)
	checkError(err)

	return nats
}
