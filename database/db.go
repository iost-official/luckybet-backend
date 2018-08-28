package database

import (
	"errors"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Database struct {
	Results   *mgo.Collection
	Rewards   *mgo.Collection
	BlockInfo *mgo.Collection
	Bets      *mgo.Collection

	todays1stRound   int
	top10            []interface{}
	lastTopQueryTime time.Time
}

func NewDatabase(session *mgo.Session) *Database {
	return &Database{
		Results:   session.DB("lucky_bet").C("results"),
		Rewards:   session.DB("lucky_bet").C("rewards"),
		BlockInfo: session.DB("lucky_bet").C("blocks"),
		Bets:      session.DB("lucky_bet").C("bets"),
	}
}

var robotAddressList = []string{"23hJissnRLwMcGFcPwyDxDfj9FaB5Z7LkY13n5TGZ2gL5", "iQ68dENWhAoAPKPa5oqvdyyhrMVJWKSAALM5wL6DaJw9", "yAw8vVunPfKrxSFCJPaxSkj4uurYNyxGcegCvyNhLkUW", "y5fZWCfJsV4kvc2RyEAxmiTg6VjX5kdsreYNnjVnshmu", "2AdmgVyJ5ystHaqcPjcjCeyt1T2A5mHDuGkyF12gmegQk", "mn8oJ62wmRQUiRhXiw9BJKu7Y9PErwdLH4fyemTF492h", "p3mA75VB39xKdj1CEbdigmxCJrS3aN3QHv3LB9Bg7btv", "yXFopFVsjiRnznaHz4DS94gGp7ESiBGZE9z6FSb4khTd", "26RnmyAXcVYyLkQroLKxAPeE5ffnXfMiP7f5uV3UqxDhc", "rgfg7KmMeQyiDzMCEGCc6qrpsfdwGC6MDSwAFPo5Knt1", "kjy5343149WyYy4PYQPYavfXPNXqV8HKCvxcu7zeLwGC", "28eRp4MEJXEoNYGqdnJ72zLdciRJgCTUjs82mchtcay3J", "tfnSXCTb3iPukhGMR3eQ9FBK8FMnwSXZevi6zt33mXhv", "26777e13XSygCdvBJfPqvvva93r6Q53qJF9gHixeZqy4B", "c8Hu3GAyQk8FM6wHiNs737ithZaNmMrKcYx2iZC84T77", "vpByWtccdakv3mGjq8YCxYkAx5mR11b39D2j4i2taNtH", "2AjmPnrRFE2EEbMtV7Hoph6RQSPaLmoK8MdsxPGDFHXwi", "22M7R5KWpe1zH2Q2qGu5iqJHNYLzsT7JtXrwtHDhChcg2", "yhuuF5qdTxEv7Lezx2t8kWa58fzevrzpVRw4UobuFcSM", "xmJdxGFQ7ErZh5EqFWp19np7YWju1Y2HETmDGmn5B4ei", "218qJFa5aFsNnQMaD5a48727d2enGXDv3NXiTG8ENUB5E", "p9sAXisoKXk2JMUwyhRND2AzRRPR7y5Xqo4EQTYU5Rhp", "xxAb6wTvEn7bFKzAKwsYQ3LYKk7PA6ZD8ARJYSyS9qqQ", "rKZnLc5ybyrzgpMVUQDN13QeApe45xDt5ZCXbxwVpPH5", "234SUfSHrEze1Z1UN3YoJTxrNuziTf6fhHradQQRkncF2", "w3d7bxCKheCuy7aKoYG9jzZybM8CERi91VxTqGRWisoD", "hyyvK2mBsCxuSLnVMYkmxaUH3FedABdv73zJsdATL9Vc", "c2gJ5v1xLpiqYqXvK8pe31ivFiuTs7qKjaMPcJqPmhwC", "h1V7zyNqSeA4ykwUHaEMvfNJsLC8cE63H7PKuUAJPFFU", "m9nNARjLZHVDRvPtgETo3hk3tPt8csdEuRWKUMebWnsJ", "sgoLawSskcdNezEHSR1cjW9WBtSGXvjfeZnZb9mkekuc", "kYfWGAdsYNhzfSBpyS6c29e5EG1AzdaSSgNWFra5yRm5", "29orQxYpoiFZmcfoms29JvmTLCzKgoPmiU2kDf6NDpEpx", "fccwWR5AvTMHKL94kUYPJfu8DHbvViaK9mupPQJJKHVc", "vU72w1hBhJwZrpYv6qfpvj4n6sWtsL6wW1xK2eLAXUxo", "j8NCUnrhsxmWhrQVQbyqFYdjKwQmBJ4VfjZEKyJ8cdUh", "caBZKqf8BfcYjydmqsHVFXZVvLry5jzTXaAT13yQoCrB", "hfNWNzLM1aDW5mNUj4YKupyiQ6oYViXA1zB9xZTS92zw", "yFXHN3YdiXsZrEvWAPLDhbSSpHxycdHhgaaujADfKSga", "nPGLceu49G1MEbTkPekWLmQbqDu9k2uJKCaSTUS524QS", "qSGdkUyrRGZo6xWxAGwSBrKGd1sDtDnoxbyyxXUBSBua", "x3NsEpH2dNSCYWKryYd2LmKXLJUXHbTcWnzWXUmRzSi9", "tBCiLVGx4DJ53wc7x8e1vCFZhkcDyVfgAET9gECYJeyX", "qPKSw34v9SeoBboNL4R1XCymXGNDoHfPuRZiDuNKrtc7", "yyyhw6oRbv4CP1XCX9xoZPu9xdifbE25xZ8U86Cvvc57", "mo6DqtLm2e6KB5Pfshjd4N9g4qnhw6Bbz2CeV2kUtFbU", "oHcqfvShKxmbyHXD7xmRjVMw5EwCSUieYbWzp1MRcSfo", "2BavaTX5fkPfu7wpK56F1aV9EekG2UGqUdKQmHvYiwiJY", "gAT86FDPyocFe79m3JHCHxMHUCpi58jPEEm2wWs2PtpV", "cCtr3P2wLn8PigtDMRfYhTXEgrbLVSbJeo5YWd8gJR88", "xA7eoJRYTRjyubJhPcZSrerF58AfBVZX9TKATPP1wr6J", "qctmYhARoYgwcTGjFgrwxReU7YYKgHCJm9XSimULgdNm", "zu75fCFZcwt26kBfQwQV5ttUn29bBVNomQMVtkf51BzV", "oPFqX3CcQ8e8gv5PSvRNDed6tapULLTLbvuEhVCYQq81", "twseF7MZHXAVjDm9kgeeiRJQrn1NmLZMa5DoTWjFoUVq", "pc6vrCCRr8PM4L8Se6KeretrNBsoQzpFkBwejWnEKy9R", "pNNibqorNipxTLJobAnp7kyNgwymkR7VPgUXZHhQoNdu", "zZ98M1S4sZ1MA9xpc3gLRy5U96F9VMzr25JmUbwGjuBv", "qverDM4JBJT8ApkEGxSQTdYTbVP7XoN1G77e7Zh5yD4n", "jKiiZmDBccMvhFchdC4jL3jPxmweiDQTSQTPSW7set4e", "zbQN3MBTgEzn97xqmEjTiQbqG8EhyMGkR7pVWrEeMzs9", "wD2ny3sNqttFgZnEVu7UWnJwLeNgwRcTPFPLQAAAAmVX", "crQwpYr6csfuPHpibY7LWBfRwxkQTpnE8ijbA4Xgq3J8", "22zgKXxBBdc85L5Jbg8X7vVfRdJGCB5xXtexSKAi8ooH8", "bYYTDeLrmiQtZshLk6b87aq7yZFR6J198GJ4ruZ818bZ", "caeR7puE7AkupTQwgbx1H14PbiiAEgb9f2Udz814tvc2", "29kJXCK3DUQpz9aJ7fpf7D9vPUUfMhXxBuBWkA3hYa2zJ", "242juxLU78QHdqWkesdDZSRxGj463UVyEtMiqCZq5EGSv", "23oSdHjdLnrMrH2jRZRZ7fEGdKugyXDsbRnLdnmwWGu9r", "22THVA7ScJcfFeFc8KCeEnUYbkft4Tk2puhjS3PaR3zEc", "oz6mMhW7dHzfeEK3jWD35hX7v8iSSsrEaC6eCxgHbvv9", "xEUMmLr3LYNN5xUNxPtisTtGfquKgRNVaJG9R2XyXp23", "23jmzULwYTseJHCjWHx546xGKpXnWrqUQioGfjwM2nYto", "jGMELeWRcFhquLN8YfUkvVxdPKjedfnjYha29YTybpja", "gNQHtgG5qL27uuT8GmxfHEHR1deNwMZGyYg96YmmotTQ", "kaV69PLsiE7a311CqcNLSQxbrzPSddk6zENz7VsEUbZb", "jGbhDg4avV4eH8Q87GUanzpDFcAewSKhWsGEZbdkrvyE", "fLyk2xZZhA9boXJUwoYWW8zxfQtiyGdShgBAkMn18DbV", "kPFCfpmxv3k7zRJxMFB6UfHxr3jHjexCBQrZgZi7Y5UQ", "zRXUubxyp4BbmmdjmYgMv2ceyVDSDhB39ePxLPdJuKgi", "mrBtmmcL3Nepywyx9mH11Cr6CVrEoAKs9ik7HNnW244K", "295PHddmC8h7n4GtAyp9knXvAecfsBRmEmN1qf7GgA2rV", "23F8Gee9BNaRAuHmhhKXoip1WqsfUkns8TcPR5fniViyG", "25eU1PAsL8AbU1NeqddeumbAKeFcVajBdpmW5REZz6tWi", "dp5wxs1xKYmQGQX38VDoUM3SL1SBmCskYgfHuM9ZouoZ", "258DESEygNhREdTr6Pviy6fRU3XHkkctKuHMBNujKcnxY", "qE954c2p4BTihaJbWc5RLMa67S8JX2cAZbxv9p1jbqfu", "vvdfSEzELRUo7KFBCRMSKWSMyr9LqVi4389uQGYZDoym", "pTqiChe4jozF625FViyzE9wV34LhwoMxSF3MTHLhaBwQ", "wN8NiTLcuCNtzSw46tUnVnYtJB6dsxkSbnbQvnwPx9Ku", "bbxumz61HFdXe3dfY2cZnTwjmcb8uZsDxwNE67TXHqJX", "xtFL5yhyNBt8egKb6SKJEGWm5hrq48bkqnmtBXqA3Mhh", "22M5QtLcfhsy7eCT37diubz4NsbGBH8MtJzQBbsYean2D", "2173azESS6CXJGgwwevm9QTP7vMoMgy2RkdSSodHHfeo2", "wC4mY2Uh1iaKb6CZVGg2wAQZxneKnKTtwHkZJmJ3yPVt", "27xqcPRLdCEqZuAqfPPMajVUoBtJ6eQ5tPMuVKpAEnb71", "24aHmprwcBpNtzRHNoATzRyXcjPufok3x6qaYMeRxMkJ1", "vT552dPv8Z9JEnq6cFfSdVR6kYNV6id3Gn8ySr1CwGgM", "fbtaNAcwNXD6LWFP3jp4Ai89QYt2LHEiaLbysZUGbmVc", "nyXbjCqEfUYFGNDb5RjvWFyCgLhG5FXrCvWwQZzwEZ23", "h6F29p52q35u4Q3LKJmyiyqXLQLnCGpNtJ7fCKoKmiuN"}

func (d *Database) Watch() {

}

type Result struct {
	Round       int
	Height      int
	LuckyNumber int
	Total       int
	Win         int
	Award       int64
}

type Reward struct {
	Round   int
	Account string
	Reward  int64
	Times   int
}

type Record struct {
	Round   int
	Account string
	Bet     int64
}

type BlockInfo struct {
	Height int
	Time   int64
}

type Bet struct {
	Address     string `json:"address"`
	LuckyNumber int    `json:"lucky_number"`
	BetAmount   int    `json:"bet_amount"`
	BetTime     int64  `json:"bet_time"`
	ClientIp    string `json:"client_ip"`
}

type top10 struct {
	Id            string `json:"id"`
	TotalWinIOST  int    `json:"total_win_iost"`
	TotalBet      int    `json:"total_bet"`
	TotalWinTimes int    `json:"total_win_times"`
	NetEarn       int    `json:"net_earn"`
}

func (d *Database) Insert(i interface{}) error {
	switch i.(type) {
	case *Result:
		d.Results.Insert(i.(*Result))
	case *Reward:
		d.Rewards.Insert(i.(*Reward))
	case *Record:
		d.Rewards.Insert(i.(*Record))
	case *BlockInfo:
		d.BlockInfo.Insert(i.(*BlockInfo))
	case *Bet:
		d.Bets.Insert(i.(*Bet))
	default:
		return errors.New("illegal type")
	}
	return nil
}

func (d *Database) QueryResult(round int) (result *Result, err error) {
	err = d.Results.Find(bson.M{"round": round}).One(&result)
	return
}

func (d *Database) QueryRewards(round int) (rewards []Reward, err error) {
	err = d.Rewards.Find(bson.M{"round": round}).All(&rewards)
	return
}

func (d *Database) QueryBlockInfo(height int) (blockInfo *BlockInfo, err error) {
	err = d.BlockInfo.Find(bson.M{"height": height}).One(&blockInfo)
	return
}

func (d *Database) QueryBet(address string, bias, length int) (bets []*Bet, err error) {
	err = d.Bets.Find(bson.M{"address": address}).Sort("-bettime").Skip(bias).Limit(length).All(&bets)
	return
}

func (d *Database) QueryBetCount(address string) int {
	n, _ := d.Bets.Find(bson.M{"address": address}).Count()
	return n
}

func (d *Database) QueryTop10(t int64) (ids []interface{}, err error) {
	if d.top10 != nil && time.Since(d.lastTopQueryTime) < 2*time.Minute {
		return d.top10, nil
	}

	queryPip := []bson.M{
		{
			"$match": bson.M{
				"round": bson.M{
					"$gte": d.todays1stRound,
				},
				"address": bson.M{
					"$nin": robotAddressList,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":           "$address",
				"totalWinIOST":  bson.M{"$sum": "$reward"},
				"totalBet":      bson.M{"$sum": "$bet"},
				"totalWinTimes": bson.M{"$sum": "$times"},
			},
		},
		{
			"$addFields": bson.M{
				"netEarn": bson.M{"$subtract": []string{"$totalWinIOST", "$totalBet"}},
			},
		},
		{
			"$sort": bson.M{
				"netEarn": -1,
			},
		},
		{
			"$limit": 10,
		},
	}

	var top10DayBetWinners []interface{}
	err = d.Rewards.Pipe(queryPip).All(&top10DayBetWinners)

	if err == nil {
		d.top10 = top10DayBetWinners
		d.lastTopQueryTime = time.Now()
	}

	return top10DayBetWinners, err
}

func (d *Database) LastBlock() *BlockInfo {
	var bi BlockInfo
	d.BlockInfo.Find(bson.M{}).Sort("-height").One(&bi)
	return &bi
}

func (d *Database) LastBet() *Result {
	var r Result
	d.Results.Find(bson.M{}).Sort("-round").One(&r)
	return &r
}
