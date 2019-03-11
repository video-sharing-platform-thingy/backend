package main

import (
	"log"
	"vspt/routes"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	// Set up routes
	router := fasthttprouter.New()
	router.POST("/upload", routes.PostUpload)

	// Listen on port 8080
	log.Println("Listening...")
	err := fasthttp.ListenAndServe(":8080", router.Handler)
	if err != nil {
		log.Fatalf("Couldn't start server: %s", err)
	}
}
