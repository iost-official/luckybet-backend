package database

import (
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
		remoteLastRound, err := Round()
		if err != nil {
			log.Fatal("watch round fatal: ", err.Error())
		}

		rw.once.Do(func() {
			rw.localLastRound, err = rw.d.QueryLastResult()
			if err != nil {
				rw.localLastRound = 0
			}
		})

		for i := rw.localLastRound + 1; i < remoteLastRound; i++ {
			r, re, err := IostResult(i)
			if err != nil {
				panic(err)
			}

			bi, err := rw.d.QueryBlockInfo(r.Height)

			if err != nil {
				r.Time = 0
			} else {
				r.Time = bi.Time
			}
			rw.d.Insert(r)

			for _, rec := range re {
				rw.d.UpdateBets(&rec, r.LuckyNumber)
				rw.d.Insert(&rec)
			}
		}

		rw.localLastRound = remoteLastRound - 1

		time.Sleep(time.Second)
	}
}

type blockWatcher struct {
	d              *Database
	localLastBlock int
}

func (bw *blockWatcher) watch() {
	for {
		remoteHeight, err := BlockChainHeight()
		if err != nil {
			panic(err)
		}

		for i := bw.localLastBlock + 1; i <= remoteHeight; i++ {
			bi, err := Block(i)
			if err != nil {
				panic(err)
			}
			bw.d.Insert(bi)
		}
		bw.localLastBlock = remoteHeight

		time.Sleep(time.Second)
	}
}
