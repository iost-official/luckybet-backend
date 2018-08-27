package database

import (
	"testing"

	"time"

	"gopkg.in/mgo.v2"
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
	}
	return &d, session
}

func TestReward(t *testing.T) {

	d, s := Init(t)
	defer s.Close()

	d.Insert(&Reward{
		Round:   0,
		Account: "abc",
		Reward:  3.14,
	})

	d.Insert(&Reward{
		Round:   0,
		Account: "def",
		Reward:  6.28,
	})

	rtn, err := d.QueryRewards(0)
	if err != nil {
		t.Fatal(err)
	}
	if rtn[0].Account != "abc" {
		t.Fatal(rtn)
	}
}

func TestResult(t *testing.T) {
	d, s := Init(t)
	defer s.Close()

	d.Insert(&Result{
		Round:       0,
		Height:      100,
		LuckyNumber: 1,
		Total:       3,
		Win:         2,
		Award:       9.42,
	})

	rtn, err := d.QueryResult(0)
	if err != nil {
		t.Fatal(err)
	}
	if rtn.Award != 9.42 {
		t.Fatal(rtn)
	}
}

func TestBlock(t *testing.T) {
	d, s := Init(t)
	defer s.Close()

	d.Insert(&BlockInfo{
		Height: 100,
		Time:   time.Now().UnixNano(),
	})

	rtn, err := d.QueryBlockInfo(100)
	if err != nil {
		t.Fatal(err)
	}
	if rtn.Height != 100 {
		t.Fatal(rtn)
	}
}
