package utils

import (
	db "SDCC/database"
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	AwsRegion              = "us-east-1"
	Replicas               = 4
	N                      = Replicas + 1
	Timeout                = 5 * time.Second
	MigrationWindowMinutes = 10
	CostRead               = 1
	CostWrite              = 2
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
		database = db.BadgerDB{Db: db.GetBadgerDb()}
	} else if parser.Database == "redis" {
		database = db.RedisDB{Db: db.ConnectToRedis()}
	} else {
		database = nil // TODO handle default
	}
	/*else if parser.Database == "go-cache" {
		database = db.GoCache{Cache: cache.New(cache.NoExpiration, 0)}
	}*/

	return Configuration{Database: database, awsRegion: parser.AwsRegion}
}
