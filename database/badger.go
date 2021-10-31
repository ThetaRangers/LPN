package database

import (
	"encoding/binary"
	"encoding/json"
	"github.com/dgraph-io/badger"
	"log"
)

type BadgerDB struct {
	Db *badger.DB
}

func (b BadgerDB) Get(key []byte) ([][]byte, uint64, error) {
	var slice [][]byte
	err := b.Db.View(func(txn *badger.Txn) error {
		var valCopy []byte
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		valCopy, err = item.ValueCopy(nil)

		err = json.Unmarshal(valCopy, &slice)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, 0, nil
		} else {
			log.Fatal(err)
			return nil, 0, err
		}
	}

	return slice[1:], binary.BigEndian.Uint64(slice[0]), nil
}

func (b BadgerDB) Put(key []byte, value [][]byte, version ...uint64) (uint64, error) {
	var versionNum uint64
	err := b.Db.Update(func(txn *badger.Txn) error {
		var err2 error
		if len(version) != 1 {
			var value [][]byte
			value, versionNum, err2 = b.Get(key)
			if err2 != nil {
				return err2
			}

			if value == nil {
				versionNum = 0
			} else {
				versionNum++
			}
		} else {
			versionNum = version[0]
		}
		entry := make([][]byte, 0)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, versionNum)
		entry = append(entry, bytes)
		entry = append(entry, value...)

		var buffer []byte
		buffer, err := json.Marshal(entry)
		if err != nil {
			log.Fatal(err)
		}
		e := badger.NewEntry(key, buffer)
		err = txn.SetEntry(e)
		return err
	})
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return versionNum, nil
}

func (b BadgerDB) Append(key []byte, value [][]byte) ([][]byte, uint64, error) {
	var versionNumber uint64
	var valCopy []byte
	var buffer []byte
	var slice [][]byte
	num := make([]byte, 8)

	err := b.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			versionNumber = 0
		} else if err != nil {
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			log.Fatal(err)
		}
		if len(valCopy) != 0 {
			err = json.Unmarshal(valCopy, &slice)
			if err != nil {
				log.Fatal(err)
			}
			versionNumber = binary.BigEndian.Uint64(slice[0])
			versionNumber++
			binary.BigEndian.PutUint64(slice[0], versionNumber)
		} else {
			binary.BigEndian.PutUint64(num, versionNumber)
			slice = append(slice, num)
		}

		slice = append(slice, value...)

		buffer, err = json.Marshal(slice)
		if err != nil {
			log.Fatal(err)
		}

		e := badger.NewEntry(key, buffer)
		err = txn.SetEntry(e)

		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return slice[1:], versionNumber, nil
}

func (b BadgerDB) Del(key []byte) error {
	err := b.Db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
