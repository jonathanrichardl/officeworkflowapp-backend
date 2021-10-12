package controller

import (
	"clean/internal/controller/models"
	"clean/internal/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (c *Controller) GetStatusOfAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	orders, err := c.order.ListOrders()
	if err != nil {
		return
	}
	response := models.BuildPayload(orders)
	for _, order := range response {
		requirements, err := c.requirements.GetRequirementsbyOrderId(order.ID)
		if err != nil {
			return
		}
		order.AddRequirements(requirements)
	}
	json.NewEncoder(w).Encode(response)

}

func (c *Controller) GetStatusOfOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	request := mux.Vars(r)
	uuid, err := uuid.FromBytes([]byte(request["id"]))
	if err != nil {
		return
	}
	order, err := c.order.GetOrder(uuid)
	response := models.BuildPayload([]*entity.Orders{order})
	requirements, err := c.requirements.GetRequirementsbyOrderId(order.ID.String())
	response[0].AddRequirements(requirements)

	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(response)

}

func (c *Controller) AddNewOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var order models.Orders
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	err = json.Unmarshal(req, &order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	deadline, _ := time.Parse("2 Jan 2006", order.Deadline)
	id, err := c.order.NewOrder(order.Title, order.Description, deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	for _, requirement := range order.Requirements {
		_, err := c.requirements.CreateRequirement(requirement.Request, requirement.ExpectedOutcome, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Request"))

		}

	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("201 - Order '%s' has been added, keep track on your order here at /orders/id=%d", order.Title, id)))

}

func (c *Controller) PostUpdateOnDelivery(w http.ResponseWriter, r *http.Request) {
	// request := mux.Vars(r)
	// id := request["id"]
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
