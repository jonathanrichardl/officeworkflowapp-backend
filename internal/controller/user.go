package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-validation-v2/internal/controller/models"
)

func (c *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	user, err := c.user.GetUserbyID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while retrieving user info: ", err.Error())
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("User Not Found"))
		return
	}
	response := models.BuildUserProfile(user)
	json.NewEncoder(w).Encode(response)

}
