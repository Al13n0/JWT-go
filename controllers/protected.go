package controllers

import (
	"log"
	"net/http"
)

func (c Controller) ProtectedEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("protected endpoint was called")
	}
}
