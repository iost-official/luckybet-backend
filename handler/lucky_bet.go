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
	Ret     int     `json:"ret"`
	Msg     string  `json:"msg"`
	Balance float64 `json:"balance"`
}

type luckyBet struct {
	Ret    int    `json:"ret"`
	Msg    string `json:"msg"`
	TxHash string `json:"tx_hash"`
}

func LuckyBet(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	address := params.ByName("address")
	betAmount := params.ByName("betAmount")
	luckyNumber := params.ByName("luckyNumber")
	privKey := params.ByName("privKey")
	gcaptcha := params.ByName("gcaptcha")

	remoteip := ctx.Request.Header.Peek("Iost_Remote_Addr")

	if !verifyGCAP(gcaptcha, string(remoteip)) {
		log.Println(ErrGreCaptcha.Error())
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrGreCaptcha.Error()})
	}

	if address == "" || betAmount == "" || privKey == "" || luckyNumber == "" {
		log.Println("GetLuckyBet nil params")
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrInvalidInput.Error()})
	}

	luckyNumberInt, err := strconv.Atoi(luckyNumber)
	if err != nil || (luckyNumberInt < 0 || luckyNumberInt > 9) {
		log.Println("GetLuckyBet invalud lucky number")
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrInvalidInput.Error()})
	}

	betAmountInt, err := strconv.Atoi(betAmount)
	if err != nil || (betAmountInt <= 0 || betAmountInt > 5) {
		log.Println("GetLuckyBet invalud bet amount")
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrInvalidInput.Error()})
	}

	if len(address) != 44 && len(address) != 45 {
		log.Println("GetLuckyBet invalid address")
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{Ret: 1, Msg: ErrInvalidInput.Error()})
	}

	balance, err := iost.BalanceByKey(address)
	if err != nil {
		log.Println("GetLuckyBet GetBalanceByKey error:", err)
	}
	if float64(betAmountInt) > balance {
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBetFail{6, ErrInsufficientBalance.Error(), balance})
	}

	// send to blockChain
	var (
		txHash        []byte
		transferIndex int
	)
	for transferIndex < 3 {
		txHash, err = iost.SendBet(address, privKey, luckyNumberInt, betAmountInt)
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
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBet{3, ErrOutOfRetryTime.Error(), ""})
	}

	txHashEncoded := common.Base58Encode(txHash)

	// check BlocakChain
	var checkIndex int
	for checkIndex < 30 {
		time.Sleep(time.Second * 2)
		if _, err := iost.GetTxnByHash(txHashEncoded); err != nil {
			//log.Printf("GetLuckyBet SendBet error: %v, retry...\n", err)
		} else {
			log.Println("GetLuckyBet blockChain Hash: ", txHashEncoded)
			break
		}
		checkIndex++
	}

	if checkIndex == 30 {
		log.Println("GetLuckyBet checkTxHash error:", ErrOutOfCheckTxHash)
		ctx.Response.Header.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(&luckyBet{4, ErrOutOfCheckTxHash.Error(), txHashEncoded})
	}
	log.Println("GetLuckyBet checkTxHash success.")

	ba := &database.Bet{
		Address:     address,
		LuckyNumber: luckyNumberInt,
		BetAmount:   betAmountInt,
		BetTime:     time.Now().Unix(),
		ClientIp:    remoteip,
	}
	D.Insert(ba)

	ctx.Response.Header.SetStatusCode(200)
	json.NewEncoder(ctx).Encode(&luckyBet{0, hex.EncodeToString(txHash), txHashEncoded})

}

func verifyGCAP(gcap, remoteip string) bool {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	args := fasthttp.AcquireArgs()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseArgs(args)

	args.Set("secret", GCAPSecretKey)
	args.Set("response", gcap)
	args.Set("remoteip", remoteip)

	//"POST", GCAPVerifyUrl, strings.NewReader(postData.Encode())

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
