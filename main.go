package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func put(key, value []byte, db badger.DB) {
	entry := make([][]byte, 0)
	entry = append(entry, value)

	var buffer []byte
	buffer, err := json.Marshal(entry)
	handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, buffer)
		err := txn.SetEntry(e)
		return err
	})
	handle(err)
}

func get(key []byte, db badger.DB) [][]byte {
	var slice [][]byte

	err := db.View(func(txn *badger.Txn) error {
		var valCopy []byte
		item, err := txn.Get(key)
		handle(err)

		valCopy, err = item.ValueCopy(nil)

		err = json.Unmarshal(valCopy, &slice)
		handle(err)

		return nil
	})

	handle(err)

	return slice
}

func append_db(key []byte, value []byte, db badger.DB) {
	err := db.Update(func(txn *badger.Txn) error {
		var valCopy []byte
		var buffer []byte

		item, err := txn.Get(key)

		valCopy, err = item.ValueCopy(nil)

		var slice [][]byte
		err = json.Unmarshal(valCopy, &slice)
		handle(err)
		slice = append(slice, value)

		buffer, err = json.Marshal(slice)
		handle(err)

		e := badger.NewEntry(key, buffer)
		err = txn.SetEntry(e)

		return err
	})
	handle(err)
}

func main() {

	//Open badger DB
	db, err := badger.Open(badger.DefaultOptions("badgerDB"))
	handle(err)

	defer func(db *badger.DB) {
		err := db.Close()
		handle(err)
	}(db)

	key := "Ciao"

	put([]byte(key), []byte("nyan"), *db)

	handle(err)

	slice := get([]byte(key), *db)

	for i := 0; i < len(slice); i++ {
		fmt.Printf("Cocci-%d: %s\n", i, slice[i])
	}

	append_db([]byte(key), []byte("cat"), *db)
	slice = get([]byte(key), *db)

	for i := 0; i < len(slice); i++ {
		fmt.Printf("Cocci-%d: %s\n", i, slice[i])
	}
}
