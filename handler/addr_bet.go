package handler

import (
	"strconv"

	"fmt"

	"encoding/json"

	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

const BetResultEachPage = 10

type betPage struct {
	BetList  []*database.Bet `json:"address_bet_list"`
	Page     int             `json:"page"`
	LastPage int             `json:"last_page"`
}

func AddressBet(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	address := params.ByName("id")
	page := params.ByName("p")
	tag := params.ByName("t")

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	pageInt, _ := strconv.Atoi(page)
	if pageInt <= 0 {
		pageInt = 1
	}

	limit := BetResultEachPage
	if tag == "t" {
		limit = 5
	}
	skip := (pageInt - 1) * limit
	top5Bet, err := D.QueryBet(address, skip, limit)
	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	totalTimes := D.QueryBetCount(address)

	var lastPage int
	if totalTimes%BetResultEachPage == 0 {
		lastPage = totalTimes / BetResultEachPage
	} else {
		lastPage = totalTimes/BetResultEachPage + 1
	}

	ctx.Response.Header.SetStatusCode(200)
	json.NewEncoder(ctx).Encode(&betPage{top5Bet, pageInt, lastPage})

}
