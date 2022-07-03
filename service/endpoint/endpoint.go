package endpoint

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"service/app"
	"service/model"
)

var orderTemplate *template.Template

func SetupRouter(r *mux.Router) {
	s := r.Methods("GET").Subrouter()
	s.HandleFunc("/order", getAllOrders)
	s.HandleFunc("/order/{id}", getOrderById)

	var err error
	orderTemplate, err = template.New("orders").ParseFiles("./resources/index.html")
	if err != nil {
		log.Fatal(err)
	}
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := app.GetService().Services.OrdersService.GetAllOrders(r.Context())
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
