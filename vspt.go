package main

import (
	"hello/routes"
	"log"

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
