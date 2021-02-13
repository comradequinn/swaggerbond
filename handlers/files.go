package handlers

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"swaggerbond/assets"
)

//SwaggerFilesDir is the relative file location that swagger files are written to and read from
var SwaggerFilesDir string

var a = map[string]string{}

func init() {
	a["vue.js"] = assets.VueJSDev
	a["bootstrap.css"] = assets.BootstrapCSS
	a["swagger-ui.css"] = assets.SwaggerUICSS
	a["swagger-ui-bundle.js"] = assets.SwaggerUIJS
	a["swaggerbond.js"] = assets.SwaggerbondJS
}

//SwaggerFile returns the requested swagger file
func SwaggerFile(w http.ResponseWriter, r *http.Request) {
	swaggerfile := path.Join(SwaggerFilesDir, path.Base(r.RequestURI))
	b, err := ioutil.ReadFile(swaggerfile)

	if err != nil {
		log.Printf("unable to read swagger file at %q\n", swaggerfile)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

//Asset returns the requested script or stylesheet
func Asset(w http.ResponseWriter, r *http.Request) {
	asset, ok := a[path.Base(r.RequestURI)]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch path.Ext(r.RequestURI) {
	case ".json":
		w.Header().Add("Content-Type", "application/json")
	case ".js":
		w.Header().Add("Content-Type", "application/javascript")
	case ".css":
		w.Header().Add("Content-Type", "text/css")
		break
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	io.WriteString(w, asset)
}
