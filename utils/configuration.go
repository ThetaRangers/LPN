package utils

import (
	db "SDCC/database"
	"encoding/json"
	"log"
	"os"
)

const (
	AwsRegion = "us-east-1"
	CostRead  = 1
	CostWrite = 2
)

var N int
var Threshold uint64
var DynamoTable string
var MigrationWindowMinutes int
var TestingServer string
var MigrationThreshold int
var TestingMode bool
var MigrationPeriodSeconds int
var RequestTimeout int
var Database db.Database
var Subnet string

func GetConfiguration() {
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
		Database               string
		ReplicationFactor      int
		OffloadingThreshold    uint64
		DynamoTable            string
		MigrationWindowMinutes int
		TestingServer          string
		MigrationThreshold     int
		TestingMode            bool
		MigrationPeriodSeconds int
		RequestTimeout         int
		DbAddress              string
		Subnet                 string
	}{}
	err = decoder.Decode(&parser)
	if err != nil {
		log.Fatal(err)
	}

	if parser.Database == "badger" {
		Database = db.BadgerDB{Db: db.GetBadgerDb()}
	} else if parser.Database == "redis" {
		Database = db.RedisDB{Db: db.ConnectToRedis(parser.DbAddress)}
	} else {
		Database = nil // TODO handle default
	}

	N = parser.ReplicationFactor
	Threshold = parser.OffloadingThreshold
	DynamoTable = parser.DynamoTable
	MigrationWindowMinutes = parser.MigrationWindowMinutes
	MigrationThreshold = parser.MigrationThreshold
	TestingServer = parser.TestingServer
	TestingMode = parser.TestingMode
	MigrationPeriodSeconds = parser.MigrationPeriodSeconds
	RequestTimeout = parser.RequestTimeout
	Subnet = parser.Subnet
}
