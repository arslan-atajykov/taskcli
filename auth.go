package main

import (
	"log"
	"net/http"
	"strings"
)

const authToken = "secret-token"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		if token != authToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("authorized")
		next.ServeHTTP(w, r) // ← это нужно вызывать обязательно
	})
}
