package handler

import (
	"log"

	"errors"

	"encoding/json"

	"strconv"

	"encoding/hex"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/iost-official/Go-IOS-Protocol/common"
	"github.com/iost-official/luckybet-backend/database"
	"github.com/iost-official/luckybet-backend/iost"
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

	lbr := luckyBetRequest{
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

type luckyBetRequest struct {
	address     string // params.ByName("address")
	betAmount   string // params.ByName("betAmount")
	luckyNumber string // params.ByName("luckyNumber")
	privKey     string // params.ByName("privKey")
	gcaptcha    string // params.ByName("gcaptcha")

	remoteip string // ctx.Request.Header.Peek("Iost_Remote_Addr")

	luckyNumberInt int
	betAmountInt   int

	txHash        []byte
	txHashEncoded string
}

func (l *luckyBetRequest) verifyGCAP() bool {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	args := fasthttp.AcquireArgs()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseArgs(args)

	args.Set("secret", GCAPSecretKey)
	args.Set("response", l.gcaptcha)
	args.Set("remoteip", l.remoteip)

	req.SetRequestURI(GCAPVerifyUrl)
	args.WriteTo(req.BodyWriter())
	req.Header.SetMethod("POST")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := gcapClient.Do(req, res)
	if err != nil {
		log.Println("verifyGCAP error:", err)
		return false
	}

	j, err := simplejson.NewJson(res.Body())
	if err != nil {
		log.Println("verifyGCAP error:", err)
		log.Println("verifyGCAP result:", string(res.Body()))
		return false
	}

	success, err := j.Get("success").Bool()
	if err != nil {
		log.Println("verifyGCAP error:", err)
		log.Println(j.EncodePretty())
		return false
	}

	return success
}

func (l *luckyBetRequest) checkArgs() bool {
	var err error
	if l.address == "" || l.betAmount == "" || l.privKey == "" || l.luckyNumber == "" {
		log.Println("GetLuckyBet nil params")
		return false
	}

	l.luckyNumberInt, err = strconv.Atoi(l.luckyNumber)
	if err != nil || (l.luckyNumberInt < 0 || l.luckyNumberInt > 9) {
		log.Println("GetLuckyBet invalud lucky number")
		return false
	}

	l.betAmountInt, err = strconv.Atoi(l.betAmount)
	if err != nil || (l.betAmountInt <= 0 || l.betAmountInt > 5) {
		log.Println("GetLuckyBet invalud bet amount")
		return false
	}

	if len(l.address) != 44 && len(l.address) != 45 {
		log.Println("GetLuckyBet invalid address")
		return false
	}
	return true
}

func (l *luckyBetRequest) checkBalance() int64 {

	balance, err := iost.BalanceByKey(l.address)
	if err != nil {
		log.Println("GetLuckyBet GetBalanceByKey error:", err)
	}
	return balance
}

func (l *luckyBetRequest) send() bool {
	var (
		txHash        []byte
		transferIndex int
	)
	for transferIndex < 3 {
		txHash, err := iost.SendBet(l.address, l.privKey, l.luckyNumberInt, l.betAmountInt)
		if err != nil {
			log.Println("GetLuckyBet SendBet error:", err)
		}
		if txHash != nil {
			break
		}
		transferIndex++
		time.Sleep(time.Second)
	}

	if transferIndex == 3 {
		log.Println("GetLuckyBet SendBet error:", ErrOutOfRetryTime)
		return false
	}

	l.txHashEncoded = common.Base58Encode(txHash)
	return true
}

func (l *luckyBetRequest) pullResult() bool {
	var checkIndex int
	for checkIndex < 30 {
		time.Sleep(time.Second * 2)
		if _, err := iost.GetTxnByHash(l.txHashEncoded); err == nil {
			log.Println("GetLuckyBet blockChain Hash: ", l.txHashEncoded)
			break
		}
		checkIndex++
	}

	if checkIndex == 30 {
		log.Println("GetLuckyBet checkTxHash error:", ErrOutOfCheckTxHash)
		return false
	}
	log.Println("GetLuckyBet checkTxHash success.")
	return true
}
