package database

import (
	"testing"
)

func TestRound(t *testing.T) {
	t.Skip("need manual")
	_, err := Round()
	if err != nil {
		t.Fatal(err)
	}
}
