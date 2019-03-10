package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func Download(ctx *fasthttp.RequestCtx) (string, error) {
	vidUuid := uuid.New().String()
	toPath, err := filepath.Abs("./uploaded/" + vidUuid + "-" + string(ctx.FormValue("name")))
	if err != nil {
		return "", err
	}

	to, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

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

	_, err = io.Copy(to, from)
	if err != nil {
		return "", err
	}
	err = to.Close()
	if err != nil {
		return "", err
	}
	err = from.Close()
	if err != nil {
		return "", err
	}

	return toPath, nil
}
