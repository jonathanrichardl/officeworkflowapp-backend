package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func (c *Controller) GetAllUncompletedOrders(w http.ResponseWriter, r *http.Request) {
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
	}
	json.NewEncoder(w).Encode(response)

}

func (c *Controller) SearchOrders(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	query := request["query"]
	orders, err := c.order.SearchOrders(query)
	if err != nil {
		c.logger.ErrorLogger.Printf("Error processing query %s : %s\n", query, err.Error())
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

	}
	json.NewEncoder(w).Encode(response)
}
func (c *Controller) AddNewOrder(w http.ResponseWriter, r *http.Request) {
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
	deadline, _ := time.Parse("2/Jan/2006 15:04:05", order.Deadline)
	id, err := c.order.NewOrder(order.Title, order.Description, deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Can't add new order into database : ", err.Error())
	}
	var wg sync.WaitGroup
	for _, requirement := range order.Requirements {
		wg.Add(1)
		go func(wg *sync.WaitGroup, requirement models.Requirements) {
			_, err := c.requirements.CreateRequirement(requirement.Request, requirement.ExpectedOutcome, id)
			if err != nil {
				c.logger.ErrorLogger.Println("Can't add requirements into database : ", err.Error())
				wg.Done()
				return
			}
			wg.Done()
		}(&wg, requirement)
	}
	wg.Wait()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("201 - Order '%s' has been added, keep track on your order here at /orders/id=%s", order.Title, id)))

}

func (c *Controller) GetStatusOfOrder(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	uuid := request["id"]
	order, err := c.order.GetOrder(uuid)

	if err != nil {
		c.logger.ErrorLogger.Printf("Error retrieving order %s from database : %s\n", uuid, err.Error())
		return
	}

	response := models.BuildPayload([]*entity.Orders{order})
	requirements, err := c.requirements.GetRequirementsbyOrderId(uuid)
	if err != nil {
		c.logger.ErrorLogger.Printf("Error retrieving requirements for order %s from database: %s\n", uuid, err.Error())
		return
	}
	fmt.Println(requirements[0].Status)
	response[0].AddRequirements(requirements)

	json.NewEncoder(w).Encode(response)

}

func (c *Controller) DeleteOrder(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) ModifyOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	var patch models.OrderPatch
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	err = json.Unmarshal(req, &patch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	orderDetail, err := c.order.GetOrder(orderID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while retrieving order : ", err.Error())
		w.Write([]byte("Internal Server Error"))
	}
	if patch.Deadline != nil {
		orderDetail.Deadline, err = time.Parse("2/Jan/2006 15:04:05", *patch.Deadline)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Deadline Datetime Format"))
			return
		}
	}
	if patch.Description != nil {
		orderDetail.Description = *patch.Description
	}

	if patch.Title != nil {
		orderDetail.Title = *patch.Title
	}

	err = c.order.UpdateOrder(orderDetail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		c.logger.ErrorLogger.Println("Error while modifying order : ", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order Modified"))

}

func (c *Controller) ModifyRequirements(w http.ResponseWriter, r *http.Request) {
	var patches models.RequirementPatch
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	err = json.Unmarshal(req, &patches)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	for _, patch := range patches.Patches {
		r, err := c.requirements.GetRequirementbyID(patch.Id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Request"))
			c.logger.ErrorLogger.Println("Error retrieving requirement : ", err.Error())
		}

		if patch.ExpectedOutcome != nil {
			r.ExpectedOutcome = *patch.ExpectedOutcome
		}

		if patch.Request != nil {
			r.Request = *patch.Request
		}

		err = c.requirements.UpdateRequirement(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error updating requirement : ", err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Requirements Modified"))

}

func (c *Controller) AddNewRequirement(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	var newRequirement models.Requirements
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	err = json.Unmarshal(req, &newRequirement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
	}
	_, err = c.order.GetOrder(orderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request, Order Does Not Exist"))
	}
	_, err = c.requirements.CreateRequirement(newRequirement.Request, newRequirement.ExpectedOutcome, orderID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected Error"))
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Requirement Added"))
}

/*
func (c *Controller) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	err := c.order.DeleteOrder(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while deleting task : ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
*/
