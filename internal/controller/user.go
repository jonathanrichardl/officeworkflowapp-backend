package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"

	"github.com/gorilla/mux"
)

func (c *Controller) GetTasks(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	tasks, err := c.requirements.GetRequirementsbyUserId(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving tasks : ", err.Error())
	}
	response := models.BuildTasks(tasks)
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
