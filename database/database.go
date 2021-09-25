package database

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"log"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Get(key []byte, db badger.DB) [][]byte {
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

func Put(key, value []byte, db badger.DB) {
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

func Append(key []byte, value []byte, db badger.DB) {
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

func Del(key []byte, db badger.DB){

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		handle(err)

		return nil
	})

	handle(err)
}