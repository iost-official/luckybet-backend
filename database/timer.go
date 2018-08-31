package database

import "github.com/hashicorp/golang-lru"

type timer struct {
	d     *Database
	cache *lru.Cache
}

func newTimer(size int, database *Database) timer {
	c, err := lru.New(size)
	if err != nil {
		panic(err)
	}
	return timer{
		cache: c,
		d:     database,
	}
}

func (t *timer) TimeOfBlock(height int) int64 {
	if v, ok := t.cache.Get(height); ok {
		return v.(int64)
	} else {
		bi, err := t.d.QueryBlockInfo(height)
		if err != nil {
			return 0
		}
		t.cache.Add(height, bi.Time)
		return bi.Time
	}
}
