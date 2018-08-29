package iost

import "github.com/iost-official/Go-IOS-Protocol/core/tx"

func BalanceByKey(address string) (int64, error) {
	return 100, nil
}

func SendBet(address, privKey string, luckyNumberInt, betAmountInt int) ([]byte, error) {
	return []byte("txhash"), nil
}

func GetTxnByHash(hash string) (*tx.Tx, error) {
	t := tx.NewTx(nil, nil, 1, 1, 100)
	return &t, nil
}
