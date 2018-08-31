package database

import (
	"time"
)

func (d *Database) Watch() {
	go d.todayWatcher.watch()
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
}

func (rw *roundWatcher) watch() {
	for {
		remoteLastRound, err := Round()
		if err != nil {
			panic(err)
		}

		for i := rw.localLastRound; i < remoteLastRound; i++ {
			r, err := IostResult(i)
			if err != nil {
				panic(err)
			}
			rw.d.Insert(r)
		}

		time.Sleep(time.Second)
	}
}
