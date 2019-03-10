package routes

import (
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"log"
	"path/filepath"
	"vspt/utils"
)

func PostUpload(ctx *fasthttp.RequestCtx) {
	log.Println("Incoming request")
	fromPath, err := utils.Download(ctx)
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
	err = utils.Transcode(fromPath, toPath)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("video/mp4")
	ctx.SendFile(toPath)
}
