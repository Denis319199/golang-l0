package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/db"
	"service/service"
	"sync/atomic"
	"time"
)

type Service struct {
	Server       *http.Server
	Database     *db.Database
	Services     *service.Services
	Validator    *validator.Validate
	Nats         stan.Conn
	shutDownFlag int32
}

var ServiceInstance *Service

func GetService() *Service {
	return ServiceInstance
}

func (service *Service) ShutdownService() {
	prev := atomic.SwapInt32(&service.shutDownFlag, 1)
	if prev == 1 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := service.Server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}

	err = service.Database.Db.Close()
	if err != nil {
		log.Println(err)
	}

	err = service.Nats.Close()
	if err != nil {
		log.Println(err)
	}
}

func (service *Service) Start() {
	go func() {
		log.Println("Service has been started")
		if err := service.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	service.ShutdownService()

	os.Exit(1)
}
