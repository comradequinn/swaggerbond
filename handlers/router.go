package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

var viewRex = regexp.MustCompile(`\/services\/.`)

//Router selects the appropriate handler for a http request
var Router http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	log.Printf("http request for %q \n", r.URL)

	switch r.Method {
	case http.MethodGet:
		switch {
		default:
			w.WriteHeader(http.StatusNotFound)
		case r.RequestURI == "/":
			Home(w, r)
		case viewRex.MatchString(r.RequestURI):
			View(w, r)
		case strings.Contains(r.RequestURI, "/services"):
			Search(w, r)
		case strings.Contains(r.RequestURI, "/assets/"):
			Asset(w, r)
		case strings.Contains(r.RequestURI, "/swagger-files/"):
			SwaggerFile(w, r)
		}
	case http.MethodPost:
		switch {
		default:
			w.WriteHeader(http.StatusNotFound)
		case strings.Contains(r.RequestURI, "/echo/"):
			Echo(w, r)
		case strings.Contains(r.RequestURI, "/capitalise/"):
			Capitalise(w, r)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
