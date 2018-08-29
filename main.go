package main

import (
	"github.com/iost-official/luckybet-backend/database"
	"github.com/iost-official/luckybet-backend/handler"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"gopkg.in/mgo.v2"
)

var router fasthttprouter.Router

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	handler.D = &database.Database{
		Results:   session.DB("test").C("results"),
		Rewards:   session.DB("test").C("rewards"),
		BlockInfo: session.DB("test").C("blocks"),
		Bets:      session.DB("test").C("bets"),

		Todays1stRound: 0,
	}

	run()
}

func run() {
	router.POST("/api/luckyBet", handler.LuckyBet)
	router.POST("/api/luckyBetBenchMark", handler.LuckyBetBenchMark)
	router.GET("/api/luckyBet/round/:id", handler.BetRound)
	router.GET("/api/luckyBet/addressBet/:id", handler.AddressBet)
	router.GET("/api/luckyBet/latestBetInfo", handler.LatestBetInfo)
	router.GET("/api/luckyBet/todayRanking", handler.TodayTop10Address)

	fasthttp.ListenAndServe(":12345", router.Handler)
}
