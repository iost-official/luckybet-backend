package nonce

import (
	"sync"

	"github.com/iost-official/luckybet-backend/database"
)

type Nonce struct {
	nonce int
	m     sync.Mutex
}

var (
	nonceInstance *Nonce
	once          sync.Once
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
