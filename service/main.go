package main

import (
	"service/app"
	_ "service/app/config"
)

func main() {
	app.ServiceInstance.Start()
}
