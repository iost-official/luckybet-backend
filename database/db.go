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

	todayWatcher
	roundWatcher
	blockWatcher

	top10            []Top10
	lastTopQueryTime time.Time
}

func NewDatabase(db *mgo.Database) *Database {
	d := &Database{
		Results:   db.C("results"),
		Rewards:   db.C("rewards"),
		BlockInfo: db.C("blocks"),
		Bets:      db.C("bets"),
	}
	d.todayWatcher = todayWatcher{
		d: d,
	}

	a, err := d.QueryLastResult()
	if err != nil {
		a = -1
	}

	d.roundWatcher = roundWatcher{
		d:              d,
		localLastRound: a,
	}

	b, err := d.QueryTopBlocks(1)
	var bb int
	if len(b) == 0 || err != nil {
		bb = 0
	} else {
		bb = b[0].Height
	}
	d.blockWatcher = blockWatcher{
		d:              d,
		localLastBlock: bb,
	}
	return d
}

var robotAddressList = []string{"23hJissnRLwMcGFcPwyDxDfj9FaB5Z7LkY13n5TGZ2gL5"}

type Result struct {
	Round       int
	Height      int `json:"number"`
	LuckyNumber int
	Total       int `json:"user_number"`
	Win         int `json:"k_number"`
	Award       int64
	Time        int64
}

//type Reward struct {
//	Round   int
//	Account string
//	Reward  int64
//	Times   int
//}

type Record struct {
	Round   int
	Account string
	Bet     int64
	Win     int64
	Nonce   int
}

type BlockInfo struct {
	Height int   `json:"height"`
	Time   int64 `json:"timestamp"`
}

type Bet struct {
	Account     string     `json:"Account"`
	LuckyNumber int        `json:"lucky_number"`
	BetAmount   int64      `json:"bet_amount"`
	BetTime     int64      `json:"bet_time"`
	ClientIp    string     `json:"client_ip"`
	Nonce       int        `json:"nonce"`
	Result      *BetResult `json:"result"`
	Status      int        `json:"status"`
}

type BetResult struct {
	Round  int   `json:"round"`
	Is_win bool  `json:"is_win"`
	Amount int64 `json:"Amount"`
}

type RoundInfo struct {
	Account  string `json:"_id"`
	Iost     int64  `json:"totalWinIost"`
	WinTimes int    `json:"totalWinTimes"`
}

type Top10 struct {
	Id            string `json:"_id"`
	TotalWinIOST  int64  `json:"totalWinIost"`
	TotalBet      int64  `json:"totalBet"`
	TotalWinTimes int    `json:"totalWinTimes"`
	NetEarn       int64  `json:"netEarn"`
}

func (d *Database) Insert(i interface{}) error {
	switch i.(type) {
	case *Result:
		d.Results.Insert(i.(*Result))
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

func (d *Database) UpdateBets(r *Record, ln int) error {
	//br := BetResult{
	//	Round:  r.Round,
	//	IsWin:  r.Win == 0,
	//	Amount: r.Win,
	//}
	res := bson.M{
		"round":  r.Round,
		"is_win": r.Win > 0,
		"amount": int(r.Win),
	}

	err := d.Bets.Update(bson.M{"nonce": r.Nonce}, bson.M{"$set": bson.M{"result": res}})

	return err
}

func (d *Database) AcceptBet(nonce int) error {
	err := d.Bets.Update(bson.M{"nonce": nonce}, bson.M{"$set": bson.M{"status": 1}})
	return err
}

func (d *Database) FailBet(nonce int) error {
	err := d.Bets.Update(bson.M{"nonce": nonce}, bson.M{"$set": bson.M{"status": -1}})
	return err
}

func (d *Database) QueryResult(round, limit int) (result []Result, err error) {
	err = d.Results.Find(bson.M{}).Limit(limit).Sort("-round").All(&result)
	return
}

func (d *Database) QueryLastResult() (last int, err error) {
	var result Result
	err = d.Results.Find(bson.M{}).Sort("-round").One(&result)
	return result.Round, err
}

func (d *Database) QueryRoundInfo(round int) (roundInfo []RoundInfo, err error) {

	query := []bson.M{
		{
			"$match": bson.M{
				"round": round,
				"win":   bson.M{"$gt": 0},
			},
		},
		{
			"$group": bson.M{
				"_id":      "$account",
				"IOST":     bson.M{"$sum": "$win"},
				"WinTimes": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{
				"IOST": -1,
			},
		},
	}
	var qr []map[string]interface{}
	err = d.Rewards.Pipe(query).All(&qr)
	roundInfo = make([]RoundInfo, len(qr))
	for i, q := range qr {
		roundInfo[i].Account = q["_id"].(string)
		roundInfo[i].Iost = q["IOST"].(int64)
		roundInfo[i].WinTimes = q["WinTimes"].(int)
	}
	return
}

func (d *Database) QueryBlockInfo(height int) (blockInfo *BlockInfo, err error) {
	err = d.BlockInfo.Find(bson.M{"height": height}).One(&blockInfo)
	return
}

func (d *Database) QueryTopBlocks(number int) (bis []BlockInfo, err error) {
	err = d.BlockInfo.Find(bson.M{}).Sort("-height").Limit(number).All(&bis)
	return
}

func (d *Database) QueryBet(address string, bias, length int) (bets []*Bet, err error) {
	err = d.Bets.Find(bson.M{"account": address}).Sort("-bettime").Skip(bias).Limit(length).All(&bets)
	return
}

func (d *Database) QueryBetCount(address string) int {
	n, _ := d.Bets.Find(bson.M{"account": address}).Count()
	return n
}

func (d *Database) QueryNonce() int {
	var bet Bet
	err := d.Bets.Find(bson.M{}).Sort("-nonce").One(&bet)
	if err != nil {
		return 0
	}
	return bet.Nonce + 1
}

func (d *Database) QueryTodays1stRound() int {
	var result Result
	err := d.Results.Find(bson.M{"time": bson.M{"$get": today().UnixNano()}}).Sort("time").One(&result)
	if err != nil {
		return -1
	}
	return result.Round
}

func (d *Database) QueryTop10(t int64) (top []Top10, err error) {
	if d.top10 != nil && time.Since(d.lastTopQueryTime) < 2*time.Minute {
		return d.top10, nil
	}

	queryPip := []bson.M{
		{
			"$match": bson.M{
				"round": bson.M{
					"$gte": d.Todays1stRound,
				},
				"account": bson.M{
					"$nin": robotAddressList,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":           "$account",
				"totalWinIOST":  bson.M{"$sum": "$win"},
				"totalBet":      bson.M{"$sum": "$bet"},
				"totalWinTimes": bson.M{"$sum": 1},
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

	var top10DayBetWinners = make([]Top10, 0)
	//err = d.Rewards.Pipe(queryPip).All(&top10DayBetWinners)

	var it []interface{}
	err = d.Rewards.Pipe(queryPip).All(&it)

	for _, m := range it {
		mr := m.(bson.M)
		top10DayBetWinners = append(top10DayBetWinners, Top10{
			Id:            mr["_id"].(string),
			TotalWinIOST:  int64(mr["totalWinIOST"].(int64)),
			TotalBet:      mr["totalBet"].(int64),
			TotalWinTimes: mr["totalWinTimes"].(int),
			NetEarn:       mr["netEarn"].(int64),
		})
	}

	d.top10 = top10DayBetWinners
	d.lastTopQueryTime = time.Now()

	return top10DayBetWinners, err
}

func (d *Database) LastBlock() *BlockInfo {
	var bi BlockInfo
	d.BlockInfo.Find(bson.M{}).Sort("-height").One(&bi)
	return &bi
}

//func (d *Database) LastBet() *Result {
//	var r Result
//	d.Results.Find(bson.M{}).Sort("-round").One(&r)
//	return &r
//}
