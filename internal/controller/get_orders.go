package controller

import (
	"clean/internal/controller/models"
	"clean/internal/entity"
	"encoding/json"
	"net/http"

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
		c.logger.ErrorLogger.Printf("Error retrieving order %s from database : %s\n", uuid, err.Error())
		return
	}

	response := models.BuildPayload([]*entity.Orders{order})
	requirements, err := c.requirements.GetRequirementsbyOrderId(uuid)
	if err != nil {
		c.logger.ErrorLogger.Printf("Error retrieving requirements for order %s from database: %s\n", uuid, err.Error())
		return
	}
	response[0].AddRequirements(requirements)

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
