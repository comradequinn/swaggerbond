package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

//Echo writes back the request payload, used for demonstration purposes
func Echo(w http.ResponseWriter, r *http.Request) {
	example(w, r, func(s string) string { return s })
}

//Capitalise writes back the request payload after sorting, used for demonstration purposes
func Capitalise(w http.ResponseWriter, r *http.Request) {
	example(w, r, func(s string) string { return strings.ToUpper(s) })
}

func example(w http.ResponseWriter, r *http.Request, action func(string) string) {
	w.Header().Set("Content-Type", "application/json")

	body := struct {
		Data string `json:"data"`
	}{}

	if e := json.NewDecoder(r.Body).Decode(&body); e != nil {
		log.Printf("example request body was invalid. %v", e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	echo := struct {
		Received string `json:"response"`
	}{Received: action(body.Data)}

	json.NewEncoder(w).Encode(echo)
}
