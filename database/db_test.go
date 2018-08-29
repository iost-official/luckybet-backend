package database

import (
	"testing"

	"time"

	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Init(t *testing.T) (*Database, *mgo.Session) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	d := Database{
		Results:   session.DB("test").C("results"),
		Rewards:   session.DB("test").C("rewards"),
		BlockInfo: session.DB("test").C("blocks"),
		Bets:      session.DB("test").C("bets"),

		Todays1stRound: 0,
	}
	return &d, session
}

func TestGenerate(t *testing.T) {
	now := int64(1535520668227538346)

	d, s := Init(t)
	defer s.Close()

	for i := 0; i < 100; i++ {
		d.Insert(&Bet{
			Account:     "player" + strconv.Itoa(i),
			LuckyNumber: i % 10,
			BetAmount:   1,
			BetTime:     now + int64(i)*int64(time.Minute),
			ClientIp:    "127.0.0.1",
		})
		d.Insert(&Record{
			Round:   i/20 + 1,
			Account: "player" + strconv.Itoa(i),
			Bet:     1,
		})
	}

	for i := 2; i < 12; i++ {
		d.Insert(&Reward{
			Round:   i / 2,
			Account: "player" + strconv.Itoa(i),
			Reward:  int64(i),
			Times:   1,
		})
	}

	for i := 1; i < 5; i++ {
		d.Insert(&Result{
			Round:       i,
			Height:      100 + 3*i,
			LuckyNumber: i,
			Total:       20,
			Win:         2,
			Award:       int64(2*i + 1),
			Time:        now + int64(15*i)*int64(time.Minute),
		})
	}

	for i := 1; i < 16; i++ {
		d.Insert(&BlockInfo{
			Height: 100 + i,
			Time:   now + int64(5*i)*int64(time.Minute),
		})
	}
}

func TestReward(t *testing.T) {

	d, s := Init(t)
	defer s.Close()

	rtn, err := d.QueryRewards(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(rtn) < 1 || rtn[0].Account != "player6" {
		t.Fatal(rtn)
	}
}

func TestResult(t *testing.T) {
	d, s := Init(t)
	defer s.Close()

	rtn, err := d.QueryResult(3, 5)
	if err != nil {
		t.Fatal(err)
	}
	if rtn[0].Height != 112 {
		t.Fatal(rtn)
	}
}

func TestBlock(t *testing.T) {
	d, s := Init(t)
	defer s.Close()

	rtn, err := d.QueryBlockInfo(110)
	if err != nil {
		t.Fatal(err)
	}
	if rtn.Height != 110 {
		t.Fatal(rtn)
	}
}

func TestClear(t *testing.T) {
	d, s := Init(t)
	defer s.Close()
	d.Rewards.RemoveAll(bson.M{})
	d.Results.RemoveAll(bson.M{})
	d.BlockInfo.RemoveAll(bson.M{})
	d.Bets.RemoveAll(bson.M{})
}
