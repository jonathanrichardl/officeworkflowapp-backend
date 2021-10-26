package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var form models.LoginForm
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Logging in: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Logging in: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	ID, role, ok, err := c.user.Login(form.Username, form.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		c.logger.ErrorLogger.Println("Error while logging in ", err.Error())
		return
	}
	if ok {
		var response models.Token
		token, err := c.generateJWT(ID, role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error generating jwt ", err.Error())
			return
		}
		response.Token = token
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Username or Password is wrong"))
		return
	}

}

func (c *Controller) generateJWT(userid string, role string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["authorization"] = role
	atClaims["exp"] = time.Now().Add(time.Minute * 15)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("SUPERSECRETPASSWORD"))
	if err != nil {
		return "", err
	}
	return token, nil

}
