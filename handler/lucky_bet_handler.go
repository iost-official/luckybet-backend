package handler

import (
	"log"
	"strconv"
	"time"

	"strings"

	"bytes"

	"github.com/bitly/go-simplejson"
	"github.com/iost-official/Go-IOS-Protocol/common"
	"github.com/iost-official/luckybet-backend/database"
	"github.com/iost-official/luckybet-backend/nonce"
	"github.com/valyala/fasthttp"
)

type luckyBetHandler struct {
	account     string // params.ByName("account")
	betAmount   string // params.ByName("betAmount")
	luckyNumber string // params.ByName("luckyNumber")
	privKey     string // params.ByName("privKey")
	gcaptcha    string // params.ByName("gcaptcha")

	remoteip string // ctx.Request.Header.Peek("Iost_Remote_Addr")

	luckyNumberInt int
	betAmountInt   int

	txHash        []byte
	txHashEncoded string

	nonce int
}

func (l *luckyBetHandler) verifyGCAP() bool {
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

func (l *luckyBetHandler) checkArgs() bool {
	var err error
	if l.account == "" || l.betAmount == "" || l.privKey == "" || l.luckyNumber == "" {
		log.Println("GetLuckyBet nil params", l.account, l.betAmount, l.privKey, l.luckyNumber)
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

	if !strings.HasPrefix(l.account, "IOST") {
		log.Println("GetLuckyBet invalid address")
		return false
	}
	return true
}

func (l *luckyBetHandler) checkBalance() int64 {

	balance, err := database.BalanceByKey(l.account)
	if err != nil {
		log.Println("GetLuckyBet GetBalanceByKey error:", err)
	}
	return balance
}

func (l *luckyBetHandler) send() bool {

	var (
		txHash        []byte
		transferIndex int
	)

	l.nonce = nonce.Instance().Get(D)

	for transferIndex < 3 {
		var err error
		txHash, err = database.SendBet(l.account, l.privKey, l.luckyNumberInt, l.betAmountInt, l.nonce)
		if err != nil {
			log.Println("GetLuckyBet SendBet error:", err)
		}
		if txHash != nil && !bytes.Equal(txHash, []byte("")) {
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

func (l *luckyBetHandler) pullResult() bool {
	var checkIndex int
	for checkIndex < 30 {
		time.Sleep(time.Second * 2)

		if _, err := database.GetTxnByHash(l.txHashEncoded); err == nil {
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

func (l *luckyBetHandler) insert() {
	ba := &database.Bet{
		Account:     l.account,
		LuckyNumber: l.luckyNumberInt,
		BetAmount:   l.betAmountInt,
		BetTime:     time.Now().Unix(),
		ClientIp:    l.remoteip,
		Nonce:       l.nonce,
	}
	D.Insert(ba)
}
