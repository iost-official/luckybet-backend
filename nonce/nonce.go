package nonce

import (
	"sync"

	"encoding/json"
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/iost-official/luckybet-backend/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type Nonce struct {
	nonce int
	m     sync.Mutex
}

var (
	nonceInstance *Nonce
	once          sync.Once
	D             *database.Database
	Client        = fasthttp.Client{}
)

func Instance() *Nonce {
	if nonceInstance == nil {
		once.Do(func() {
			nonceInstance = newNonce()
		})
	}
	return nonceInstance
}

func newNonce() *Nonce {
	return &Nonce{
		nonce: -1,
	}
}

func (n *Nonce) Get(d *database.Database) int {
	n.m.Lock()
	defer n.m.Unlock()
	if n.nonce < 0 {
		n.nonce = d.QueryNonce()
	}
	rtn := n.nonce
	n.nonce++
	return rtn
}

type NonceRes struct {
	Nonce int `json:"nonce"`
}

func Handler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.Response.SetStatusCode(200)
	nr := NonceRes{
		Nonce: Instance().Get(D),
	}
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(ctx).Encode(nr)
	if err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func ReqNonce(url string) int {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	err := Client.Do(req, res)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	j, err := simplejson.NewJson(res.Body())
	if err != nil {
		fmt.Println(err)
		return -1
	}
	rtn, err := j.Get("nonce").Int()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return rtn
}
