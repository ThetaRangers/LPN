package utils

import (
	db "SDCC/database"
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
)

type Configuration struct {
	Database  db.Database
	awsRegion string
}

func GetConfiguration() Configuration {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	parser := struct {
		Database  string
		AwsRegion string
	}{}
	err = decoder.Decode(&parser)
	if err != nil {
		log.Fatal(err)
	}

	var database db.Database
	if parser.Database == "badger" {
		badgerDB, err := badger.Open(badger.DefaultOptions("badgerDB"))
		if err != nil {
			log.Fatal(err)
		}
		database = db.BadgerDB{Db: badgerDB}
	} else if parser.Database == "go-cache" {
		database = db.GoCache{Cache: cache.New(cache.NoExpiration, 0)}
	} else if parser.Database == "redis" {
		database = db.RedisDB{Db: db.ConnectToRedis()}
	} else {
		database = nil // TODO handle default
	}

	return Configuration{Database: database, awsRegion: parser.AwsRegion}
}
