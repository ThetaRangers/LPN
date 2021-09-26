package database

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"log"
	"sync"
)

var once sync.Once

type BadgerDB struct {
	db *badger.DB
}

var instance *BadgerDB

func getInstance() *BadgerDB {
	once.Do(func() {
		database, err := badger.Open(badger.DefaultOptions("badgerDB"))
		if err != nil {
			log.Fatal(err)
		}
		instance = &BadgerDB{db: database}
	})
	return instance
}

func (b *BadgerDB) getDB() badger.DB {
	return *b.db
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b BadgerDB) Get(key []byte) [][]byte {
	var slice [][]byte
	db := getInstance().getDB()
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

func (b BadgerDB) Put(key, value []byte) {
	entry := make([][]byte, 0)
	entry = append(entry, value)

	var buffer []byte
	buffer, err := json.Marshal(entry)
	handle(err)

	db := getInstance().getDB()
	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, buffer)
		err := txn.SetEntry(e)
		return err
	})
	handle(err)
}

func (b BadgerDB) Append(key, value []byte) {
	db := getInstance().getDB()
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

func (b BadgerDB) Del(key []byte) {
	db := getInstance().getDB()
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		handle(err)

		return nil
	})

	handle(err)
}
