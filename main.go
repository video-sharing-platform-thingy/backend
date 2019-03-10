package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"vspt/routes"
	"vspt/utils"
)

func main() {
	router := fasthttprouter.New()
	router.POST("/upload", routes.PostUpload)

	log.Println("Listening...")
	err := fasthttp.ListenAndServe(":8080", router.Handler)
	utils.Check(err)
}
