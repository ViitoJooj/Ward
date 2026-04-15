package database

import "github.com/cockroachdb/pebble"

var DB *pebble.DB

func Conn() {
	var err error
	DB, err = pebble.Open("./data.db", &pebble.Options{
		MemTableSize: 64 << 20,
	})
	if err != nil {
		panic(err)
	}
}
