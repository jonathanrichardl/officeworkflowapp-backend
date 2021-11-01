package controller

import (
	"net/http"
	"order-validation-v2/internal/usecase/orders"
	"order-validation-v2/internal/usecase/requirements"
	"order-validation-v2/internal/usecase/submissions"
	"order-validation-v2/internal/usecase/tasks"
	"order-validation-v2/internal/usecase/user"
	"order-validation-v2/pkg/logger"
	"os"

	"github.com/gorilla/mux"
)

type Controller struct {
	router       *mux.Router
	order        orders.UseCase
	user         user.UseCase
	task         tasks.UseCase
	submissions  submissions.UseCase
	requirements requirements.UseCase
	logger       *logger.LoggerInstance
}

func NewController(o orders.UseCase, u user.UseCase, r requirements.UseCase, t tasks.UseCase, s submissions.UseCase, l *logger.LoggerInstance) *Controller {
	router := mux.NewRouter().StrictSlash(true)
	controller := &Controller{router: router, order: o, user: u, requirements: r, task: t, submissions: s, logger: l}
	return controller
}

func (c *Controller) RegisterHandler() {
	login := c.router.PathPrefix("/login").Subrouter()
	login.HandleFunc("/", c.Login).Methods("POST")

	userapp := c.router.PathPrefix("/orders").Subrouter()
	userapp.Use(c.validateUserJWT)
	userapp.HandleFunc("/", c.GetTasks).Methods("GET")
	userapp.HandleFunc("/submission", c.PostUpdateOnTask).Methods("POST")

	admin := c.router.PathPrefix("/admin").Subrouter()
	admin.Use(c.validateAdminJWT)
	admin.HandleFunc("/orders", c.GetStatusOfAllOrders).Methods("GET")
	admin.HandleFunc("/orders", c.AddNewOrder).Methods("POST")
	admin.HandleFunc("/orders/id={id}", c.GetStatusOfOrder).Methods("GET")
	admin.HandleFunc("/orders/id={id}", c.DeleteOrder).Methods("DELETE")
	admin.HandleFunc("/orders/id={id}", c.ModifyRequirements).Methods("PATCH")
	admin.HandleFunc("/orders/search:{query}", c.SearchOrders).Methods("GET")
	admin.HandleFunc("/user", c.NewUser).Methods("POST")
	admin.HandleFunc("/tasks", c.AddNewTask).Methods("POST")

}

func (c *Controller) Start() {
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, c.router)
}
