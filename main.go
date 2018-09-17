package main

import (
	"fmt"

	"os"

	"github.com/iost-official/luckybet-backend/database"
	"github.com/iost-official/luckybet-backend/handler"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"gopkg.in/mgo.v2"
)

var router fasthttprouter.Router

//var contractAddress = "3uhVHYjKoK6XXwmG4H2TKPX5Fd8BvfZdKZU9U8a6EMMp"

func main() {

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//database.Contract = "Contract" + contractAddress

	database.Contract = os.Args[1]

	database.LocalIServer = "http://" + os.Args[2]

	err = session.DB("lucky_bet").C("bets").EnsureIndexKey("account", "nonce", "bettime")
	if err != nil {
		fmt.Println(err)
	}
	err = session.DB("lucky_bet").C("results").EnsureIndexKey("time", "round")
	if err != nil {
		fmt.Println(err)
	}
	err = session.DB("lucky_bet").C("rewards").EnsureIndexKey("round", "account")
	if err != nil {
		fmt.Println(err)
	}
	err = session.DB("lucky_bet").C("blocks").EnsureIndexKey("height")
	if err != nil {
		fmt.Println(err)
	}

	handler.D = database.NewDatabase(session.DB("lucky_bet"))

	if len(os.Args) > 3 {
		go handler.D.Watch()
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
	router.GET("/api/luckyBetBlockInfo", handler.BetInfo)

	fasthttp.ListenAndServe(":12345", router.Handler)
}
