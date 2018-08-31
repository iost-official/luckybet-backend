package database

import (
	"fmt"
	"testing"
)

func TestRound(t *testing.T) {
	//t.Skip("need manual")
	i, err := Round()
	t.Log(i)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValue(t *testing.T) {
	//t.Skip("need manual")
	j, err := value("result0")
	t.Log(j)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIostResult(t *testing.T) {
	r, re, rw, err := IostResult(0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r, re, rw)
}
