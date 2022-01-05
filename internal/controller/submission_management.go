package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"sync"

	"github.com/gorilla/mux"
)

func (c *Controller) ReviewSubmission(w http.ResponseWriter, r *http.Request) {
	adminID := fmt.Sprintf("%v", r.Context().Value(ctxKey{}))
	submissionID := mux.Vars(r)["id"]
	var reviewForm models.ReviewForm
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request: ", err.Error())
		return
	}
	err = json.Unmarshal(req, &reviewForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request: ", err.Error())
		return
	}
	submission, err := c.submissions.GetSubmission(submissionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid SubmissionID"))
		c.logger.ErrorLogger.Println("Invalid Request: ", err.Error())
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go c.processReviewForm(adminID, submission.TaskID, &wg, reviewForm.Approved, reviewForm.ForwardTo, reviewForm.Message)
	wg.Wait()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Review Success"))

}
