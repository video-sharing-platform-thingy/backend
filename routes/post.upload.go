package routes

import (
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func PostUpload(ctx *fasthttp.RequestCtx) {
	log.Println("Incoming request")
	fromPath, err := download(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	uuid := uuid.New().String()
	toPath, err := filepath.Abs("./transcoded/" + uuid + ".mp4")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	transcode(fromPath, toPath)

	ctx.SetContentType("video/mp4")
	ctx.SendFile(toPath)
}
