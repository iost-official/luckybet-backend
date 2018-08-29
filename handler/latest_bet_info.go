package handler

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type latestBetInfo struct {
	Height int        `json:"block_height"`
	Time   int64      `json:"block_time"`
	List   []*betInfo `json:"latest_bet_list"`
}

type betInfo struct {
	Round  int   `json:"round"`
	Height int   `json:"BlockHeight"`
	Total  int   `json:"TotalUserNumber"`
	Win    int   `json:"WinUserNumber"`
	Award  int64 `json:"TotalRewards"`
	Time   int64 `json:"win_time"`
}

func LatestBetInfo(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	bi, err := D.QueryBlockInfo(D.LastBlock().Height)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	rtn := latestBetInfo{
		Height: bi.Height,
		Time:   bi.Time,
	}

	rtn.List = make([]*betInfo, 0)
	last5, err := D.QueryResult(0, 5)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	for _, r := range last5 {
		bi := betInfo{
			Round:  r.Round,
			Height: r.Height,
			Total:  r.Total,
			Win:    r.Win,
			Award:  r.Award,
			Time:   r.Time,
		}
		rtn.List = append(rtn.List, &bi)
	}

	err = json.NewEncoder(ctx).Encode(rtn)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetStatusCode(200)
}
