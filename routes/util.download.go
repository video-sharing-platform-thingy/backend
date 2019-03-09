package routes

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func download(ctx *fasthttp.RequestCtx) (string, error) {
	uuid := uuid.New().String()
	toPath, err := filepath.Abs("./uploaded/" + uuid + "-" + string(ctx.FormValue("name")))
	if err != nil {
		return "", err
	}
	to, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}

	fheader, err := ctx.FormFile("file")
	if err != nil {
		return "", err
	}
	from, err := fheader.Open()
	if err != nil {
		return "", err
	}

	io.Copy(to, from)
	to.Close()
	from.Close()

	return toPath, nil
}
