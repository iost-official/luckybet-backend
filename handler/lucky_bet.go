package handler

import (
	"log"

	"errors"

	"encoding/json"

	"encoding/hex"
	"time"

	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

const (
	GCAPVerifyUrl = "https://www.google.com/recaptcha/api/siteverify"
	GCAPSecretKey = "6Lc1vF8UAAAAAGv1XihAK4XygBMn3UobipWMqBym"
)

var (
	ErrGreCaptcha          = errors.New("reCAPTCHA check failed")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrOutOfRetryTime      = errors.New("out of retry time")
	ErrOutOfCheckTxHash    = errors.New("out of check txHash retry time")
)

var gcapClient fasthttp.Client

type luckyBetFail struct {
	Ret     int    `json:"ret"`
	Msg     string `json:"msg"`
	Balance int64  `json:"balance"`
}

type luckyBet struct {
	Ret    int    `json:"ret"`
	Msg    string `json:"msg"`
	TxHash string `json:"tx_hash"`
}

func LuckyBet(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {

	lbr := luckyBetHandler{
		address:     params.ByName("address"),
		betAmount:   params.ByName("betAmount"),
		luckyNumber: params.ByName("luckyNumber"),
		privKey:     params.ByName("privKey"),
		gcaptcha:    params.ByName("gcaptcha"),

		remoteip: string(ctx.Request.Header.Peek("Iost_Remote_Addr")),
	}
	address := params.ByName("address")

	if !lbr.verifyGCAP() {
		log.Println(ErrGreCaptcha.Error())
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrGreCaptcha.Error()})
		return
	}

	if !lbr.checkArgs() {
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrInvalidInput.Error()})
		return
	}

	balance := lbr.checkBalance()
	if balance < int64(lbr.betAmountInt) {
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{6, ErrInsufficientBalance.Error(), balance})
		return
	}

	if !lbr.send() {
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBet{3, ErrOutOfRetryTime.Error(), ""})
		return
	}

	if !lbr.pullResult() {
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBet{4, ErrOutOfCheckTxHash.Error(), lbr.txHashEncoded})
	}

	ba := &database.Bet{
		Address:     address,
		LuckyNumber: lbr.luckyNumberInt,
		BetAmount:   lbr.betAmountInt,
		BetTime:     time.Now().Unix(),
		ClientIp:    lbr.remoteip,
	}
	D.Insert(ba)

	ctx.Response.Header.SetStatusCode(200)
	json.NewEncoder(ctx).Encode(&luckyBet{0, hex.EncodeToString(lbr.txHash), lbr.txHashEncoded})

}
