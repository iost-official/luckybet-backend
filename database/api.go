package database

import (
	"strconv"

	"encoding/json"

	"fmt"

	"errors"

	"log"

	"time"

	"github.com/bitly/go-simplejson"
	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/tx"
	"github.com/iost-official/go-iost/crypto"
	"github.com/valyala/fasthttp"
)

var (
	LocalIServer = "http://52.192.65.220:30001/"
	Client       = fasthttp.Client{
		MaxIdleConnDuration: 10 * time.Minute,
		MaxConnsPerHost:     6000,
	}
	Contract = "Contract" + "AC5V12562T7XB74A8gBe3cjfwWDbJheLWjzyY8VL6JPK"
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

type RawTxReq struct {
	Data []byte `json:"data,omitempty"`
}

func SendBet(address, privKey string, luckyNumberInt int, betAmountInt int64, nonce int) ([]byte, error) {
	act := tx.NewAction(Contract, "bet", fmt.Sprintf(`["%v",%d,%d,%d]`, address, luckyNumberInt, betAmountInt, nonce))

	te := time.Now().Add(50 * time.Second).UnixNano()

	t := tx.NewTx([]*tx.Action{&act}, nil, 100000, 1, te)
	a, err := account.NewAccount(common.Base58Decode(privKey), crypto.Ed25519)
	if err != nil {
		return nil, err
	}

	t, err = tx.SignTx(t, a)
	if err != nil {
		return nil, err
	}

	b := RawTxReq{
		Data: t.Encode(),
	}
	j, err := json.Marshal(b)
	var res *simplejson.Json
	res, err = post(LocalIServer+"sendRawTx", j)
	if err != nil {
		return nil, err
	}
	p, err := res.EncodePretty()
	if err != nil {
		return nil, err
	}
	log.Println(string(p))

	return t.Hash(), nil
}

func GetTxnByHash(hash string) (*tx.Tx, error) {

	// /getTxByHash/{hash}
	j, err := get(LocalIServer + "getTxByHash/" + hash)
	if err != nil {
		return nil, err
	}

	var t tx.Tx
	if _, ok := j.CheckGet("hash"); !ok {
		return nil, fmt.Errorf("not found")
	}
	b, err := j.Encode()
	json.Unmarshal(b, &t)
	return &t, nil
}

func BlockChainHeight() (int, error) {
	j, err := get(LocalIServer + "getHeight")
	if err != nil {
		return 0, err
	}

	s, err := j.Get("height").String()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(s)
}

func Block(height int) (*BlockInfo, error) {
	j, err := get(LocalIServer + fmt.Sprintf("getBlockByNum/%v/0", height))
	if err != nil {
		return nil, err
	}

	hs, err := j.Get("head").Get("number").String()
	if err != nil {
		return nil, fmt.Errorf("get block: %v", err)
	}
	h, _ := strconv.Atoi(hs)

	ts, err := j.Get("head").Get("time").String()
	if err != nil {
		return nil, fmt.Errorf("get block: %v", err)
	}
	t, _ := strconv.ParseInt(ts, 10, 64)

	bi := &BlockInfo{
		Height: h,
		Time:   t * 3,
	}
	return bi, nil
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

func IostResult(round int) (*Result, []Record, error) {
	j, err := value("result" + strconv.Itoa(round))
	if err != nil {
		return nil, nil, err
	}
	s, _ := j.String()
	if err != nil {
		return nil, nil, err
	}

	buf := []byte(s[1:])

	jbuf, err := simplejson.NewJson(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("parse buf: %v", err)
	}

	var res Result
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshal: %v", err)
	}

	res.Round = round
	res.LuckyNumber = res.Height % 10

	awardstr, err := jbuf.Get("total_coins").Get("number").String()
	if err != nil {
		return nil, nil, err
	}

	res.Award, err = strconv.ParseInt(awardstr, 10, 64)
	if err != nil {
		return nil, nil, err
	}

	records := make([]Record, 0)

	a := jbuf.Get("records").MustArray([]interface{}{})
	for _, m := range a {
		rec := Record{
			Round: round,
		}
		rec.Account = m.(map[string]interface{})["account"].(string)
		rec.Bet, err = m.(map[string]interface{})["coins"].(json.Number).Int64()
		if err != nil {
			return nil, nil, err
		}
		n, err := m.(map[string]interface{})["nonce"].(json.Number).Int64()
		if err != nil {
			return nil, nil, err
		}
		win, ok := m.(map[string]interface{})["reward"]
		if !ok {
			rec.Win = 0
		} else {
			w, ok := win.(string)
			if !ok {
				log.Println("invalid reward ", win)
				return nil, nil, errors.New("invalid reward")
			}
			w2, err := strconv.ParseInt(w, 10, 64)
			if err != nil {
				log.Println("invalid reward ", w)
				return nil, nil, errors.New("invalid reward")
			}
			rec.Win = w2
		}

		rec.Nonce = int(n)
		records = append(records, rec)
	}

	return &res, records, err
}

func value(key string) (*simplejson.Json, error) {
	j, err := get(LocalIServer + "getState/" + Contract + "-" + key)
	//fmt.Println("api212", LocalIServer+"getState/"+Contract+"-"+key)
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
		return nil, fmt.Errorf("get: %v", err)
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
		return nil, fmt.Errorf("post: %v", err)
	}

	return simplejson.NewJson(res.Body())
}
