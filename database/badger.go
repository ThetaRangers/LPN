package database

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
)

type BadgerDB struct {
	Db *badger.DB
}

func (b BadgerDB) Get(key []byte) ([][]byte, error) {
	var slice [][]byte
	err := b.Db.View(func(txn *badger.Txn) error {
		var valCopy []byte
		item, err2 := txn.Get(key)
		if err2 != nil {
			return err2
		}

		valCopy, err2 = item.ValueCopy(nil)

		err2 = json.Unmarshal(valCopy, &slice)
		if err2 != nil {
			return err2
		}

		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return slice, nil
}

func (b BadgerDB) Put(key []byte, value [][]byte) error {
	err := b.Db.Update(func(txn *badger.Txn) error {
		var buffer []byte
		buffer, err2 := json.Marshal(value)
		if err2 != nil {
			return err2
		}
		e := badger.NewEntry(key, buffer)
		err2 = txn.SetEntry(e)
		return err2
	})
	if err != nil {
		return err
	}

	return nil
}

func (b BadgerDB) Append(key []byte, value [][]byte) ([][]byte, error) {
	var versionNumber uint64
	var valCopy []byte
	var buffer []byte
	var slice [][]byte

	err := b.Db.Update(func(txn *badger.Txn) error {
		item, err2 := txn.Get(key)
		if err2 == badger.ErrKeyNotFound {
			versionNumber = 0
		} else if err2 != nil {
			return err2
		}

		valCopy, err2 = item.ValueCopy(nil)
		if err2 != nil {
			return err2
		}

		if len(valCopy) != 0 {
			err2 = json.Unmarshal(valCopy, &slice)
			if err2 != nil {
				return err2
			}
		}

		slice = append(slice, value...)

		buffer, err2 = json.Marshal(slice)
		if err2 != nil {
			return err2
		}

		e := badger.NewEntry(key, buffer)
		err2 = txn.SetEntry(e)

		return err2
	})
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (b BadgerDB) Del(key []byte) error {
	err := b.Db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b BadgerDB) DeleteExcept(keys []string) error {
	err := b.Db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			if !Contains(keys, string(k)) {
				err := txn.Delete(k)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func GetBadgerDb() *badger.DB {
	badgerDB, err := badger.Open(badger.DefaultOptions("badgerDB"))
	if err != nil {
		panic(err)
	}
	return badgerDB
}
