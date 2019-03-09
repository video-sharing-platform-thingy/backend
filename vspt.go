package main

import (
	"log"
	"vspt/routes"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.POST("/upload", routes.PostUpload)

	log.Println("Listening...")
	err := fasthttp.ListenAndServe(":8080", router.Handler)
	check(err)
}
