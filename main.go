package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"github.com/video-sharing-platform-thingy/backend/routes"
	"github.com/video-sharing-platform-thingy/backend/utils"
)

func main() {
	router := fasthttprouter.New()
	router.POST("/upload", routes.PostUpload)

	log.Println("Listening...")
	err := fasthttp.ListenAndServe(":8080", router.Handler)
	utils.Check(err)
}
