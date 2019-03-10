package routes

import (
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func PostUpload(ctx *fasthttp.RequestCtx) {
	log.Println("Incoming request")
	fromPath, err := Download(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	vidUuid := uuid.New().String()
	toPath, err := filepath.Abs("./transcoded/" + vidUuid + ".mp4")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = Transcode(fromPath, toPath)
	if err != nil{
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("video/mp4")
	ctx.SendFile(toPath)
}
