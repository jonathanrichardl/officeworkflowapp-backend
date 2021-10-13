package controller

import (
	"clean/internal/usecase/orders"
	"clean/internal/usecase/requirements"
	"clean/pkg/logger"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Controller struct {
	router       *mux.Router
	order        orders.UseCase
	requirements requirements.UseCase
	logger       *logger.LoggerInstance
}

func NewController(o orders.UseCase, r requirements.UseCase, l *logger.LoggerInstance) *Controller {
	router := mux.NewRouter().StrictSlash(true)
	controller := &Controller{router: router, order: o, requirements: r, logger: l}
	return controller
}

func (c *Controller) RegisterHandler() {
	c.router.HandleFunc("/", c.Index).Methods("GET")
	c.router.HandleFunc("/orders", c.GetStatusOfAllOrders).Methods("GET")
	c.router.HandleFunc("/orders", c.AddNewOrder).Methods("POST")
	c.router.HandleFunc("/orders/id={id}", c.GetStatusOfOrder).Methods("GET")
	c.router.HandleFunc("/orders/id={id}", c.PostUpdateOnDelivery).Methods("POST")
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func (c *Controller) Start() {
	http.ListenAndServe(":"+os.Getenv("PORT"), c.router)
}
