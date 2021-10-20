package controller

import (
	"net/http"
)

func (c *Controller) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if storedCookie, _ := r.Cookie("Authentication"); storedCookie != nil {
			next.ServeHTTP(w, r)
		}
	})
}
