package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
)

func (c *Controller) GetTasks(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	tasks, err := c.task.GetTasksofUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while getting tasks: ", err.Error())
		return
	}
	fmt.Println("OK")
	if tasks == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) PostUpdateOnTask(w http.ResponseWriter, r *http.Request) {
	var submission models.Submission
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Submitting: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &submission)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while unmarshalling submission: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	id, err := c.submissions.NewSubmission(submission.Message, submission.Images, submission.TaskID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while adding submission: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Submission has been accepted, id = %s", id)))

}
