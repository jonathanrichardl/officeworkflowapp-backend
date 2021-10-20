package controller

import (
	"clean/internal/controller/models"
	"clean/internal/usecase/orders"
	"clean/internal/usecase/requirements"
	"clean/internal/usecase/user"
	"clean/pkg/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Controller struct {
	router       *mux.Router
	order        orders.UseCase
	user         user.UseCase
	requirements requirements.UseCase
	logger       *logger.LoggerInstance
}

func NewController(o orders.UseCase, u user.UseCase, r requirements.UseCase, l *logger.LoggerInstance) *Controller {
	router := mux.NewRouter().StrictSlash(true)
	controller := &Controller{router: router, order: o, user: u, requirements: r, logger: l}
	return controller
}

func (c *Controller) RegisterHandler() {
	c.router.HandleFunc("/login", c.Login).Methods("POST")
	c.router.HandleFunc("/orders", c.GetStatusOfAllOrders).Methods("GET")
	c.router.HandleFunc("/orders", c.AddNewOrder).Methods("POST")
	c.router.HandleFunc("/orders/id={id}", c.GetStatusOfOrder).Methods("GET")
	c.router.HandleFunc("/orders/id={id}", c.PostUpdateOnDelivery).Methods("POST")
	c.router.HandleFunc("/orders/id={id}", c.DeleteOrder).Methods("DELETE")
	c.router.HandleFunc("/orders/id={id}", c.ModifyRequirements).Methods("PATCH")
	c.router.HandleFunc("/orders/search:{query}", c.SearchOrders).Methods("GET")
	c.router.Use(c.loggingMiddleware)

}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var form models.LoginForm
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	err = json.Unmarshal(req, &form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	ID, ok, err := c.user.Login(form.Username, form.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while logging in ", err.Error())
	}
	if ok {
		cookie := &http.Cookie{}
		cookie.Name = "Authentication"
		cookie.Value = ID
		cookie.Expires = time.Now().Add(1 * time.Hour)
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login Successful"))

	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Username or Password is wrong"))
	}

}

func (c *Controller) Start() {
	http.ListenAndServe(":8080", c.router)
}
