package main

import (
	"bufio"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	//Open badger DB
	db, err := badger.Open(badger.DefaultOptions("badgerDB"))
	handle(err)

	defer func(db *badger.DB) {
		err := db.Close()
		handle(err)
	}(db)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter key: ")
	key, _ := reader.ReadString('\n')

	fmt.Print("Enter value: ")
	value, _ := reader.ReadString('\n')

	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value))
		err := txn.SetEntry(e)
		return err
	})

	var _, valCopy []byte

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		handle(err)

		err = item.Value(func(val []byte) error {
			return nil
		})

		valCopy, err = item.ValueCopy(nil)
		return nil
	})

	fmt.Printf("The value is: %s\n", valCopy)

}
