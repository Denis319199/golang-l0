package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"service/app"
	"service/db"
	"service/endpoint"
	"service/service"
	"time"
)

func init() {
	app.ServiceInstance = &app.Service{}

	setupConfiguration()
	app.ServiceInstance.Validator = newValidator()
	app.ServiceInstance.Database = db.NewDatabase()
	app.ServiceInstance.Services = service.NewServices(app.ServiceInstance.Database)
	app.ServiceInstance.Nats = newNatsStreaming()
	app.ServiceInstance.Server = newServer()

	log.Println("Service has been successfully configured")
}

func newValidator() *validator.Validate {
	return validator.New()
}

func newServer() *http.Server {
	viper.SetDefault("service.port", ":8080")

	port := viper.GetString("service.port")

	if port[0:1] != ":" {
		port = ":" + port
	}

	r := mux.NewRouter()
	endpoint.SetupRouter(r)

	return &http.Server{
		Addr:         "0.0.0.0" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
}

func setupConfiguration() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./resources")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
}
