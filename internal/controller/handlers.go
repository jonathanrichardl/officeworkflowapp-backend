package controller

import (
	"clean/internal/controller/models"
	"clean/internal/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (c *Controller) GetStatusOfAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	orders, err := c.order.ListOrders()
	if err != nil {
		c.logger.ErrorLogger.Println("Error retrieving orders from database: ", err.Error())
		return
	}
	response := models.BuildPayload(orders)
	for _, order := range response {
		requirements, err := c.requirements.GetRequirementsbyOrderId(order.ID)
		if err != nil {
			c.logger.ErrorLogger.Println("Error retrieving requirements from database: ", err.Error())
			return
		}

		order.AddRequirements(requirements)
		fmt.Println(order)
	}
	json.NewEncoder(w).Encode(response)

}

func (c *Controller) GetStatusOfOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	request := mux.Vars(r)
	uuid := request["id"]
	order, err := c.order.GetOrder(uuid)

	if err != nil {
		c.logger.ErrorLogger.Printf("Error retrieving order %s from database table orders : %s\n", uuid, err.Error())
		return
	}

	response := models.BuildPayload([]*entity.Orders{order})
	requirements, err := c.requirements.GetRequirementsbyOrderId(uuid)
	if err != nil {
		c.logger.ErrorLogger.Printf("Error retrieving requirements for order %s from database table requirements: %s\n", uuid, err.Error())
		return
	}
	response[0].AddRequirements(requirements)

	json.NewEncoder(w).Encode(response)

}

func (c *Controller) AddNewOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprint(w, "orders")
	var order models.Orders
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request")
		return
	}
	err = json.Unmarshal(req, &order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request, Can't unmarshal :", err.Error())
		return
	}
	deadline, _ := time.Parse("2 Jan 2006", order.Deadline)
	id, err := c.order.NewOrder(order.Title, order.Description, deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Can't add new order into database : ", err.Error())
	}
	for _, requirement := range order.Requirements {
		_, err := c.requirements.CreateRequirement(requirement.Request, requirement.ExpectedOutcome, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Request"))
			c.logger.ErrorLogger.Println("Can't add requirements into database : ", err.Error())
			return

		}

	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("201 - Order '%s' has been added, keep track on your order here at /orders/id=%s", order.Title, id)))

}

func (c *Controller) PostUpdateOnDelivery(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	_ = request["id"]
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var form models.ProgressForm
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
	for _, fufillment := range form.Fufillments {
		r, err := c.requirements.GetRequirementbyID(fufillment.Requirementid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Request"))
		}
		if r.ExpectedOutcome == fufillment.Outcome {
			r.SetStatus(true)
			c.requirements.UpdateRequirement(r)
			continue
		}
		r.SetStatus(false)
		c.requirements.UpdateRequirement(r)
	}

}

func (c *Controller) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	request := mux.Vars(r)
	uuid := request["id"]

	requirements, err := c.requirements.GetRequirementsbyOrderId(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(requirements) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, requirement := range requirements {
		err = c.requirements.DeleteRequirement(requirement.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	err = c.order.DeleteOrder(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
