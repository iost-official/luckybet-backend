package database

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Database struct {
	Results   *mgo.Collection
	Rewards   *mgo.Collection
	BlockInfo *mgo.Collection
}

func NewDatabase(session *mgo.Session) *Database {
	return &Database{
		Results:   session.DB("lucky_bet").C("results"),
		Rewards:   session.DB("lucky_bet").C("rewards"),
		BlockInfo: session.DB("lucky_bet").C("blocks"),
	}
}

func (d *Database) Watch() {

}

type Result struct {
	Round       int
	Height      int
	LuckyNumber int
	Total       int
	Win         int
	Award       float64
}

type Reward struct {
	Round   int
	Account string
	Reward  float64
}

type BlockInfo struct {
	Height int
	Time   int64
}

func (d *Database) Insert(i interface{}) error {
	switch i.(type) {
	case *Result:
		d.Results.Insert(i.(*Result))
	case *Reward:
		d.Rewards.Insert(i.(*Reward))
	case *BlockInfo:
		d.BlockInfo.Insert(i.(*BlockInfo))
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
