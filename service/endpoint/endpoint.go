package endpoint

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"service/app"
	"service/model"
	"strconv"
)

var orderTemplate *template.Template

func SetupRouter(r *mux.Router) {
	s := r.Methods("GET").Subrouter()
	s.HandleFunc("/order", getOrders)
	s.HandleFunc("/order/{id}", getOrderById)

	var err error
	orderTemplate, err = template.New("orders").ParseFiles("./resources/index.html")
	if err != nil {
		log.Fatal(err)
	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var (
		page, size int
		err        error
	)

	page, err = strconv.Atoi(values.Get("page"))
	if err != nil {
		processError(w, err.Error(), 400)
	}

	size, err = strconv.Atoi(values.Get("size"))
	if err != nil {
		processError(w, err.Error(), 400)
	}

	orders, err := app.GetService().Services.OrdersService.GetOrders(page, size, r.Context())
	if err != nil {
		processError(w, err.Error(), 400)
		return
	}

	err = orderTemplate.ExecuteTemplate(w, "T", orders)
	if err != nil {
		processError(w, err.Error(), 500)
	}
}

func getOrderById(w http.ResponseWriter, r *http.Request) {
	orderUid := mux.Vars(r)["id"]

	order, err := app.GetService().Services.OrdersService.GetOrderById(r.Context(), orderUid)
	if err != nil {
		processError(w, err.Error(), 400)
		return
	}

	err = orderTemplate.ExecuteTemplate(w, "T", []*model.Order{order})
	if err != nil {
		processError(w, err.Error(), 500)
	}
}
