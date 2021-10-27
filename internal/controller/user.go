package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"

	"github.com/gorilla/mux"
)

func (c *Controller) GetTasks(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	tasks, err := c.requirements.GetRequirementsbyUserId(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving tasks : ", err.Error())
		return
	}
	if tasks == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	response := models.BuildTasks(tasks)
	m := make(map[string]entity.Orders)
	for _, task := range response {
		if order, ok := m[task.OrderID]; ok {
			task.OrderDescription = order.Description
			task.OrderDeadline = order.Deadline.Format("2006-01-02 15:04:05")
			task.OrderTitle = order.Title
			continue
		}
		order, _ := c.order.GetOrder(task.OrderID)
		task.OrderDescription = order.Description
		task.OrderDeadline = order.Deadline.Format("2006-01-02 15:04:05")
		task.OrderTitle = order.Title
		m[task.OrderID] = *order
	}

	json.NewEncoder(w).Encode(response)
}

func (c *Controller) PostUpdateOnTask(w http.ResponseWriter, r *http.Request) {

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
