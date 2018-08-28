package iost

import "github.com/iost-official/Go-IOS-Protocol/core/tx"

func BalanceByKey(address string) (int64, error) {

}

func SendBet(address, privKey string, luckyNumberInt, betAmountInt int) ([]byte, error) {

}

func GetTxnByHash(hash string) (*tx.Tx, error) {

}
