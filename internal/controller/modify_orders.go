package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"time"

	"github.com/gorilla/mux"
)

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
	deadline, _ := time.Parse("2 Jan 2006", order.Deadline)
	id, err := c.order.NewOrder(order.Title, order.Description, deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Can't add new order into database : ", err.Error())
	}
	for _, requirement := range order.Requirements {
		_, err := c.requirements.CreateRequirement(requirement.Request, requirement.ExpectedOutcome, id, &requirement.UserID)
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

func (c *Controller) ModifyRequirements(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	_ = request["id"]
	w.Header().Set("Content-Type", "application/json")
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

		if patch.UserID != nil {
			r.AssignUser(*patch.UserID)
		}
		if patch.ExpectedOutcome != nil {
			r.ExpectedOutcome = *patch.ExpectedOutcome
		}
		err = c.requirements.UpdateRequirement(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error updating requirement : ", err.Error())

		}
	}

}
