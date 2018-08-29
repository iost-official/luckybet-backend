package handler

import (
	"strconv"

	"sort"

	"encoding/json"

	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type roundInfo struct {
	Id       string `json:"_id"`
	Iost     int64  `json:"totalWinIost"`
	WinTimes int    `json:"totalWinTimes"`
}
type roundInfos []roundInfo

func (r roundInfos) Len() int {
	return len(r)
}
func (r roundInfos) Less(i, j int) bool {
	return r[i].Iost > r[j].Iost
}
func (r roundInfos) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func BetRound(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	s := params.ByName("id")
	round, err := strconv.Atoi(s)
	if err != nil {
		ctx.Response.SetStatusCode(406)
	}

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	r, err := D.QueryRewards(round)

	var infos = make([]roundInfo, 0)

	for _, reward := range r {
		infos = append(infos, roundInfo{
			Id:       reward.Account,
			Iost:     reward.Reward,
			WinTimes: 1,
		})
	}

	sort.Sort(roundInfos(infos))

	err = json.NewEncoder(ctx).Encode(infos)
	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
