package handler

import (
	"time"

	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func TodayTop10Address(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.Header.SetStatusCode(200)

	t, err := D.QueryTop10(time.Now().UnixNano())

	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	err = json.NewEncoder(ctx).Encode(t)
	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
