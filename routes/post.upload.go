package routes

import (
	"log"
	"path/filepath"
	"vspt/utils"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func PostUpload(ctx *fasthttp.RequestCtx) {
	log.Println("Incoming request")
	fromPath, err := utils.Download(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	vidUUID := uuid.New().String()
	toPath, err := filepath.Abs("./transcoded/" + vidUUID + ".mp4")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = utils.Transcode(fromPath, toPath)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("video/mp4")
	ctx.SendFile(toPath)
}
