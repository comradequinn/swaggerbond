package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"swaggerbond/assets"
	"swaggerbond/handlers"
	"swaggerbond/index"
	"time"
)

func main() {
	log.SetPrefix("swaggerbond ")

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic in main: %s", err)
		}
	}()

	p := flag.Int("p", 8080, "the (p)ort to listen on")
	d := flag.String("d", "swagger-files", "the (d)irectory to read swagger files from")
	i := flag.Int("i", 5, "the (i)interval in seconds to poll for swagger file directory changes")
	s := flag.Bool("s", false, "whether to generate an (s)ample swagger files on startup")

	flag.Parse()

	handlers.SwaggerFilesDir = *d

	if _, err := os.Stat(*d); os.IsNotExist(err) {
		log.Printf("creating swagger file directory at %q\n", *d)
		os.Mkdir(*d, 0777)
	}

	if *s {
		assets.DemoFiles(*d, fmt.Sprintf(":%v", *p))
	}

	index.Start(time.Duration(*i)*time.Second, *d)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", *p), handlers.Router); err != nil {
			log.Fatalf("unable to start http server. %v", err)
		}
	}()

	log.Printf("listening on port %v", *p)

	select {}
}
