package handler

import (
	"encoding/json"
	"fmt"

	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type betInfos struct {
	Ret     int                  `json:"ret"`
	Top6Blk []database.BlockInfo `json:"top_6_blk"`
}

func BetInfo(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	top6, err := D.QueryTopBlocks(6)
	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	bi := betInfos{
		Ret:     0,
		Top6Blk: top6,
	}

	json.NewEncoder(ctx).Encode(bi)
	ctx.SetStatusCode(200)
}
