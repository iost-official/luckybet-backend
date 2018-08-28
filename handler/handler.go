package handler

import (
	"github.com/iost-official/luckybet-backend/database"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

var D *database.Database

func Init() {
}
