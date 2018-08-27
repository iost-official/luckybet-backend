package handler

import (
	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

var D *database.Database

func Init() {

}

func LuckyBetBenchMark(*fasthttp.RequestCtx, fasthttprouter.Params) {}

func AddressBet(*fasthttp.RequestCtx, fasthttprouter.Params)        {}
func TodayTop10Address(*fasthttp.RequestCtx, fasthttprouter.Params) {}
