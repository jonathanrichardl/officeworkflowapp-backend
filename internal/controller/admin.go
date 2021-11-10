package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"
	"time"

	"github.com/gorilla/mux"
)

func (c *Controller) NewUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.NewUser
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request: ", err.Error())
		return
	}
	err = json.Unmarshal(req, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request, Can't unmarshal :", err.Error())
		return
	}
	id, err := c.user.CreateUser(newUser.Username, newUser.Email, newUser.Password, newUser.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while creating new user: ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User %s has been added with id %s\n", newUser.Username, id)))
}

func (c *Controller) GetStatusOfAllOrders(w http.ResponseWriter, r *http.Request) {
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
	deadline, _ := time.Parse("2 Jan 2006", order.Deadline)
	id, err := c.order.NewOrder(order.Title, order.Description, deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Can't add new order into database : ", err.Error())
	}
	for _, requirement := range order.Requirements {
		fmt.Println(requirement.UserID)
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

func (c *Controller) GetStatusOfOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

func (c *Controller) AddNewTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.NewTask
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}
	id, err := c.task.CreateTask(newTask.RequirementID, newTask.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error creating new task: ", err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Task %s has been created for user %s\n", id, newTask.UserID)))
}

func (c *Controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.user.ListUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving all user: ", err.Error())
		return
	}
	var response []models.RetrievedUser
	for _, user := range users {
		response = append(response, models.BuildUserProfile(user))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
