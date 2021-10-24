package controller

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func (c *Controller) validateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("JWT")
		if authorization == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authorization, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("SUPERSECRETPASSWORD"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id := fmt.Sprintf("%v", claims["user_id"])
		r.Header.Set("ID", id)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)

	})
}
