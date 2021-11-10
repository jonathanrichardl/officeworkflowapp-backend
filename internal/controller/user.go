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
	tasks, err := c.task.GetTasksofUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while getting tasks: ", err.Error())
		return
	}
	if tasks == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No Tasks Present"))
		return
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetSubmission(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	id := request["id"]
	submission, err := c.submissions.GetSubmission(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while getting submissions: ", err.Error())
		return
	}
	response := models.BuildSubmissionPayload(submission)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (c *Controller) PostSubmission(w http.ResponseWriter, r *http.Request) {
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
	image := models.DecodeSubmissionPayload(submission)
	id, err := c.submissions.NewSubmission(submission.Message, image, submission.TaskID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while adding submission: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	task, err := c.task.Get(submission.TaskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while updating task: ", err.Error())
		return
	}
	task.Status = 1
	err = c.task.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while updating task: ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Submission has been accepted, id = %s", id)))

}

func (c *Controller) UpdateSubmission(w http.ResponseWriter, r *http.Request) {

}
