package handler

import (
	"strconv"

	"encoding/json"

	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func BetRound(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	s := params.ByName("id")
	round, err := strconv.Atoi(s)
	if err != nil {
		ctx.Response.SetStatusCode(406)
	}

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	r, err := D.QueryRoundInfo(round)

	err = json.NewEncoder(ctx).Encode(r)
	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
