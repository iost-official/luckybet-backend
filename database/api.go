package database

import (
	"strconv"

	"encoding/json"

	"fmt"

	"time"

	"github.com/bitly/go-simplejson"
	"github.com/iost-official/Go-IOS-Protocol/account"
	"github.com/iost-official/Go-IOS-Protocol/common"
	"github.com/iost-official/Go-IOS-Protocol/core/tx"
	"github.com/iost-official/Go-IOS-Protocol/rpc"
	"github.com/valyala/fasthttp"
)

var (
	LocalIServer = "http://localhost:30301/"
	Client       = fasthttp.Client{}
	Contract     = "Contract2qzmSGZA6cup8BNhzhbss5sP6ivy4EQQHjaVuB2q2GTw"
)

func BalanceByKey(address string) (int64, error) {

	j, err := get(LocalIServer + "getBalance/" + address + "/0")
	if err != nil {
		return 0, err
	}

	str, err := j.Get("balance").String()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(str, 10, 64)
}

func SendBet(address, privKey string, luckyNumberInt, betAmountInt int) ([]byte, error) {
	act := tx.NewAction(Contract, "bet", fmt.Sprintf(`["%v",%d,%d]`, address, luckyNumberInt, betAmountInt))
	t := tx.NewTx([]*tx.Action{&act}, nil, 10000, 1, time.Now().UnixNano()+10*time.Second.Nanoseconds())
	a, err := account.NewAccount(common.Base58Decode(privKey))
	if err != nil {
		return nil, err
	}

	t, err = tx.SignTx(t, a)
	if err != nil {
		return nil, err
	}

	b := rpc.RawTxReq{
		Data: t.Encode(),
	}
	j, err := json.Marshal(b)
	_, err = post(LocalIServer+"/sendRawTx", j)
	if err != nil {
		return nil, err
	}

	return t.Hash(), nil
}

func GetTxnByHash(hash string) (*tx.Tx, error) {

	// /getTxByHash/{hash}
	j, err := get(LocalIServer + "getTxByHash/" + hash)
	if err != nil {
		return nil, err
	}

	var t tx.Tx
	b, err := j.Encode()
	json.Unmarshal(b, &t)
	return &t, nil
}

func Round() (int, error) {
	s, err := value("round")
	if err != nil {
		return 0, err
	}
	ss, err := s.String()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(ss[1:])
}

func IostResult(round int) (*Result, error) {
	j, err := value("result" + strconv.Itoa(round))
	if err != nil {
		return nil, err
	}
	js, err := j.Bytes()
	if err != nil {
		return nil, err
	}
	var res Result
	err = json.Unmarshal(js, &res)
	return &res, err
}

func value(key string) (*simplejson.Json, error) {
	j, err := get(LocalIServer + "getState/" + Contract + "-" + key)
	if err != nil {
		fmt.Println("get err :", err)
		return nil, err
	}

	return j.Get("value"), nil
}

func get(url string) (*simplejson.Json, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	err := Client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return simplejson.NewJson(res.Body())
}

func post(url string, body []byte) (*simplejson.Json, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(body)

	err := Client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return simplejson.NewJson(res.Body())
}
