package iost

import (
	"strconv"

	"encoding/json"

	"github.com/bitly/go-simplejson"
	"github.com/iost-official/Go-IOS-Protocol/core/tx"
	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
)

var (
	LocalIServer = "localhost:30301/"
	Client       = fasthttp.Client{}
	Contract     = "ContractFPcQWT3io6DSekcoY72waon3racgwbdPp5ULScC1W9A5"
)

func init() {

}

func BalanceByKey(address string) (int64, error) {

	j, err := get(LocalIServer + "getBalance/" + address + "/0")
	if err != nil {
		return 0, err
	}

	return j.Get("balance").Int64()
}

func SendBet(address, privKey string, luckyNumberInt, betAmountInt int) ([]byte, error) {
	return []byte("txhash"), nil
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

func Result(round int) (*database.Result, error) {
	j, err := value("result" + strconv.Itoa(round))
	if err != nil {
		return nil, err
	}
	js, err := j.Bytes()
	if err != nil {
		return nil, err
	}
	var res database.Result
	err = json.Unmarshal(js, &res)
	return res, err
}

func value(key string) (*simplejson.Json, error) {
	j, err := get(LocalIServer + "getState/" + Contract + "-" + key)
	if err != nil {
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

func post(url string, body map[string]string) (*simplejson.Json, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	j, _ := json.Marshal(body)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(j)

	err := Client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return simplejson.NewJson(res.Body())
}
