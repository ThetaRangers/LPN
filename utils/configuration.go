package utils

import (
	db "SDCC/database"
	"encoding/json"
	"log"
	"os"
	"time"
)

var N int
var Threshold uint64
var DynamoTable string
var MigrationWindowTime time.Duration
var TestingServer string
var MigrationThreshold uint64
var TestingMode bool
var MigrationPeriodTime time.Duration
var RequestTimeout time.Duration
var Database db.Database
var Subnet string
var AwsRegion string
var CostRead uint64
var CostWrite uint64
var MaxKey int

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
		Database                   string
		ReplicationFactor          int
		OffloadingThreshold        uint64
		DynamoTable                string
		MigrationWindowMinutes     uint
		TestingServer              string
		MigrationThreshold         uint64
		TestingMode                bool
		MigrationPeriodSeconds     uint
		RequestTimeoutMilliseconds uint
		DbAddress                  string
		Subnet                     string
		AwsRegion                  string
		CostRead                   uint64
		CostWrite                  uint64
		MaxKey                     int
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
	MigrationWindowTime = time.Minute * time.Duration(parser.MigrationWindowMinutes)
	MigrationThreshold = parser.MigrationThreshold
	TestingServer = parser.TestingServer
	TestingMode = parser.TestingMode
	MigrationPeriodTime = time.Second * time.Duration(parser.MigrationPeriodSeconds)
	RequestTimeout = time.Millisecond * time.Duration(parser.RequestTimeoutMilliseconds)
	Subnet = parser.Subnet
	AwsRegion = parser.AwsRegion
	CostRead = parser.CostRead
	CostWrite = parser.CostWrite
	MaxKey = parser.MaxKey
}
