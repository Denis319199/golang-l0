package endpoint

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"service/app"
	"service/service"
	"testing"
)

func SetupRouterTest(r *mux.Router) {
	s := r.Methods("GET").Subrouter()
	s.HandleFunc("/order", getAllOrders)
	s.HandleFunc("/order/{id}", getOrderById)

	var err error
	orderTemplate, err = template.New("orders").ParseFiles("./../resources/index.html")
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetOrderById(t *testing.T) {
	app.ServiceInstance = &app.Service{Services: &service.Services{}}
	app.ServiceInstance.Services.OrdersService = service.CreateOrdersServiceWithCache(false)

	r := mux.NewRouter()
	SetupRouterTest(r)
	server := httptest.NewServer(r)
	defer server.Close()

	request, _ := http.NewRequest("GET", server.URL+"/order/1", nil)
	response, _ := http.DefaultClient.Do(request)

	if response.StatusCode != 200 {
		t.Error("Wrong error code")
	}
}

func TestGetNotExistingOrderById(t *testing.T) {
	app.ServiceInstance = &app.Service{Services: &service.Services{}}
	app.ServiceInstance.Services.OrdersService = service.CreateOrdersServiceWithCache(true)

	r := mux.NewRouter()
	SetupRouterTest(r)
	server := httptest.NewServer(r)
	defer server.Close()

	request, _ := http.NewRequest("GET", server.URL+"/order/1", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != 400 {
		t.Error("Wrong error code")
	}
}

func TestGetAllOrders(t *testing.T) {
	app.ServiceInstance = &app.Service{Services: &service.Services{}}
	app.ServiceInstance.Services.OrdersService = service.CreateOrdersServiceWithCache(false)

	r := mux.NewRouter()
	SetupRouterTest(r)
	server := httptest.NewServer(r)
	defer server.Close()

	request, _ := http.NewRequest("GET", server.URL+"/order", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		t.Error("Wrong error code")
	}
}
