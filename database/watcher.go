package database

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func (d *Database) Watch() {
	go d.todayWatcher.watch()
	go d.blockWatcher.watch()

	time.Sleep(3 * time.Second)
	go d.roundWatcher.watch()
}

type todayWatcher struct {
	d              *Database
	Todays1stRound int
	today          int64
}

func (tw *todayWatcher) watch() {
	for {
		t := today().UnixNano()
		if tw.today < t {
			tw.today = t
		}
		tw.Todays1stRound = tw.d.QueryTodays1stRound()
		time.Sleep(time.Minute)
	}

}

func today() time.Time {
	dateStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	return t
}

type roundWatcher struct {
	d              *Database
	localLastRound int
	once           sync.Once
}

func (rw *roundWatcher) watch() {
	for {
		time.Sleep(time.Second)

		remoteLastRound, err := Round()
		if err != nil {
			log.Println("watch round Error: ", err.Error())
		}

		rw.once.Do(func() {
			rw.localLastRound, err = rw.d.QueryLastResult()
			if err != nil {
				rw.localLastRound = 0
			}
		})

		for i := rw.localLastRound + 1; i < remoteLastRound; {
			r, re, err := IostResult(i)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			bi, err := rw.d.QueryBlockInfo(r.Height)

			if err != nil {
				fmt.Println("query time err", r.Height)
				time.Sleep(time.Second)
				continue
			}

			r.Time = bi.Time
			rw.d.Insert(r)

			for _, rec := range re {
				rw.d.UpdateBets(&rec, r.LuckyNumber)
				rw.d.Insert(&rec)
			}
			i++
		}

		rw.localLastRound = remoteLastRound - 1

	}
}

type blockWatcher struct {
	d              *Database
	localLastBlock int
}

func (bw *blockWatcher) watch() {
Outer:
	for {
		time.Sleep(time.Second)

		remoteHeight, err := BlockChainHeight()
		if err != nil {
			log.Println("watch failed", err)
			continue Outer
		}

		for i := bw.localLastBlock + 1; i <= remoteHeight; i++ {
			bi, err := Block(i)
			if err != nil {
				log.Println("pull block info failed", err)
				continue Outer
			}
			bw.d.Insert(bi)
		}
		bw.localLastBlock = remoteHeight

	}
}
