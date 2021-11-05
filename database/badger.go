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

func (b BadgerDB) Put(key []byte, value [][]byte) (uint64, error) {
	var versionNum uint64
	var slice [][]byte

	err := b.Db.Update(func(txn *badger.Txn) error {
		var err2 error
		item, err2 := txn.Get(key)
		if err2 == badger.ErrKeyNotFound {
			versionNum = 0
		} else if err2 != nil {
			return err2
		} else {
			valCopy, err2 := item.ValueCopy(nil)
			if err2 != nil {
				return err2
			}
			err2 = json.Unmarshal(valCopy, &slice)
			if err2 != nil {
				return err2
			}

			versionNum = binary.BigEndian.Uint64(slice[0])
			versionNum++
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

func (b BadgerDB) Replicate(key []byte, value [][]byte, version uint64) error {
	var versionNum uint64
	var slice [][]byte

	err := b.Db.Update(func(txn *badger.Txn) error {
		var err2 error
		item, err2 := txn.Get(key)
		if err2 != nil && err2 != badger.ErrKeyNotFound {
			return err2
		} else if err2 != badger.ErrKeyNotFound {
			valCopy, err2 := item.ValueCopy(nil)
			if err2 != nil {
				return err2
			}

			err2 = json.Unmarshal(valCopy, &slice)
			if err2 != nil {
				return err2
			}

			versionNum = binary.BigEndian.Uint64(slice[0])

			// Replica does this
			if versionNum > version {
				// Replica has the newer value
				return nil
			}
		}

		versionNum = version

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
		return err
	}

	return nil
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
			return err
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

func (b BadgerDB) Migrate(key []byte) ([][]byte, uint64, error) {
	var slice [][]byte

	err := b.Db.Update(func(txn *badger.Txn) error {
		item, err2 := txn.Get(key)
		if err2 == badger.ErrKeyNotFound {
			return nil
		} else if err2 != nil {
			return err2
		}

		valCopy, err2 := item.ValueCopy(nil)

		err2 = json.Unmarshal(valCopy, &slice)
		if err2 != nil {
			return err2
		}

		err2 = txn.Delete(key)
		if err2 != nil {
			return err2
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return slice[1:], binary.BigEndian.Uint64(slice[0]), nil
}

func GetBadgerDb() *badger.DB {
	badgerDB, err := badger.Open(badger.DefaultOptions("badgerDB"))
	if err != nil {
		panic(err)
	}
	return badgerDB
}
