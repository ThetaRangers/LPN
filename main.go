package main

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("answer"), []byte("42"))
		err := txn.SetEntry(e)
		return err
	})

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			log.Fatal(err)
		}

		var _, valCopy []byte
		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.

			// Accessing val here is valid.
			fmt.Printf("The answer is: %s\n", val)

			// Copying or parsing val is valid.
			valCopy = append([]byte{}, val...)
			return nil
		})

		// You must copy it to use it outside item.Value(...).
		fmt.Printf("The answer is: %s\n", valCopy)

		// Alternatively, you could also use item.ValueCopy().
		valCopy, err = item.ValueCopy(nil)
		fmt.Printf("The answer is: %s\n", valCopy)

		return nil
	})

}
