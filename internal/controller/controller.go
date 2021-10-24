package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/usecase/orders"
	"order-validation-v2/internal/usecase/requirements"
	"order-validation-v2/internal/usecase/user"
	"order-validation-v2/pkg/logger"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	login := c.router.PathPrefix("/login").Subrouter()
	login.HandleFunc("/", c.Login).Methods("POST")
	orders := c.router.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("/", c.GetStatusOfAllOrders).Methods("GET")
	orders.HandleFunc("/", c.AddNewOrder).Methods("POST")
	orders.HandleFunc("/id={id}", c.GetStatusOfOrder).Methods("GET")
	orders.HandleFunc("/id={id}", c.PostUpdateOnDelivery).Methods("POST")
	orders.HandleFunc("/id={id}", c.DeleteOrder).Methods("DELETE")
	orders.HandleFunc("/id={id}", c.ModifyRequirements).Methods("PATCH")
	orders.HandleFunc("/search:{query}", c.SearchOrders).Methods("GET")
	orders.Use(c.validateJWT)

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
		var response models.Token
		token, err := c.generateJWT(ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error generating jwt ", err.Error())
		}
		response.Token = token
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Username or Password is wrong"))
	}

}

func (c *Controller) Start() {
	http.ListenAndServe(":8080", c.router)
}

func (c *Controller) generateJWT(userid string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("SUPERSECRETPASSWORD"))
	if err != nil {
		return "", err
	}
	return token, nil

}
