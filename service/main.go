package main

import (
	"service/app"
	_ "service/app/config"
)

func main() {
	// import config package to configure service
	app.ServiceInstance.Start()
}
