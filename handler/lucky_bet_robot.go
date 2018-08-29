package handler

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func LuckyBetBenchMark(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	lbr := luckyBetHandler{
		account:     params.ByName("address"),
		betAmount:   params.ByName("betAmount"),
		luckyNumber: params.ByName("luckyNumber"),
		privKey:     params.ByName("privKey"),
		gcaptcha:    params.ByName("gcaptcha"),

		remoteip: string(ctx.Request.Header.Peek("Iost_Remote_Addr")),
	}
	address := params.ByName("address")

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.SetStatusCode(200)

	if !lbr.checkArgs() {
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrInvalidInput.Error()})
		return
	}

	balance := lbr.checkBalance()
	if balance < int64(lbr.betAmountInt) {
		json.NewEncoder(ctx).Encode(&luckyBetFail{6, ErrInsufficientBalance.Error(), balance})
		return
	}

	if !lbr.send() {
		json.NewEncoder(ctx).Encode(&luckyBet{3, ErrOutOfRetryTime.Error(), ""})
		return
	}

	if !lbr.pullResult() {
		json.NewEncoder(ctx).Encode(&luckyBet{4, ErrOutOfCheckTxHash.Error(), lbr.txHashEncoded})
	}

	ba := &database.Bet{
		Account:     address,
		LuckyNumber: lbr.luckyNumberInt,
		BetAmount:   lbr.betAmountInt,
		BetTime:     time.Now().Unix(),
		ClientIp:    lbr.remoteip,
	}
	D.Insert(ba)

	json.NewEncoder(ctx).Encode(&luckyBet{0, hex.EncodeToString(lbr.txHash), lbr.txHashEncoded})
}
