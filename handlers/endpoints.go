package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"swaggerbond/assets"
	"swaggerbond/index"
)

//Home writes the initial HTML for Swaggerbond
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, assets.HomeHTML)
}

//Search writes the result of performing the specified in JSON
func Search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("search")
	results := index.Search(term)

	log.Printf("search for %q returned %v results ", term, len(results))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

//View writes the SwaggerUI response for the requested service
func View(w http.ResponseWriter, r *http.Request) {
	slug := strings.ToLower(path.Base(r.RequestURI))

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, strings.Replace(assets.SwaggerUIHTML, "{# service-slug #}", slug, 1))
}
