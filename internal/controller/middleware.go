package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type ctxKey struct{}

func (c *Controller) validateUserJWT(next http.Handler) http.Handler {
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
			c.logger.ErrorLogger.Println("Error processing JWT: ", err.Error())
			return
		}
		id := fmt.Sprintf("%v", claims["user_id"])
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
func (c *Controller) validateAdminJWT(next http.Handler) http.Handler {
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
			c.logger.ErrorLogger.Println("Error processing JWT: ", err.Error())
			return
		}
		id := fmt.Sprintf("%v", claims["user_id"])
		role := fmt.Sprintf("%v", claims["authorization"])
		if role != "Admin" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
