package handler

import (
	"encoding/hex"
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func LuckyBetBenchMark(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	params := ctx.PostArgs()
	lbr := luckyBetHandler{
		account:     string(params.Peek("address")),
		betAmount:   string(params.Peek("betAmount")),
		luckyNumber: string(params.Peek("luckyNumber")),
		privKey:     string(params.Peek("privKey")),
		gcaptcha:    string(params.Peek("gcaptcha")),

		remoteip: string(ctx.Request.Header.Peek("Iost_Remote_Addr")),
	}

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
		return
	}

	lbr.insert()

	json.NewEncoder(ctx).Encode(&luckyBet{0, hex.EncodeToString(lbr.txHash), lbr.txHashEncoded})
}
