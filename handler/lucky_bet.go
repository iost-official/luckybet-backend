package handler

import (
	"errors"
	"log"

	"encoding/json"

	"encoding/hex"

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

func LuckyBet(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
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

	if !lbr.verifyGCAP() {
		log.Println(ErrGreCaptcha.Error())
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrGreCaptcha.Error()})
		return
	}

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

	lbr.insert()

	json.NewEncoder(ctx).Encode(&luckyBet{0, hex.EncodeToString(lbr.txHash), lbr.txHashEncoded})

}
