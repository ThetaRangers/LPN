package database

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisDB struct {
	Db *redis.Client
}

func (r RedisDB) Get(key []byte) ([][]byte, uint64, error) {
	ctx := context.Background()
	var slice [][]byte

	val, err := r.Db.Get(ctx, string(key)).Bytes()
	if err == redis.Nil {
		return nil, 0, nil
	} else if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(val, &slice)
	if err != nil {
		return nil, 0, err
	}

	return slice[1:], binary.BigEndian.Uint64(slice[0]), nil
}

func (r RedisDB) Put(key []byte, value [][]byte, version ...uint64) (uint64, error) {
	ctx := context.Background()
	var versionNum uint64

	fmt.Println("Starting transaction")

	_, err := r.Db.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		if len(version) != 1 {
			var value [][]byte
			var err2 error

			value, versionNum, err2 = r.Get(key)
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

		buffer, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		err = r.Db.Set(ctx, string(key), buffer, 0).Err()
		fmt.Println("Set data")
		if err != nil {
			return err
		}

		return nil
	})
	fmt.Println("Ending transaction")
	if err != nil {
		return 0, err
	}

	return versionNum, nil
}

func (r RedisDB) Append(key []byte, value [][]byte) ([][]byte, uint64, error) {
	ctx := context.Background()
	var slice [][]byte
	var versionNumber uint64
	var num = make([]byte, 8)

	_, err := r.Db.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		val, err := pipe.Get(ctx, string(key)).Bytes()
		if err != nil {
			return err
		}
		if len(val) != 0 {
			err = json.Unmarshal(val, &slice)
			if err != nil {
				log.Fatal(err)
			}
			versionNumber = binary.BigEndian.Uint64(slice[0])

			versionNumber++
			binary.BigEndian.PutUint64(slice[0], versionNumber)
		} else {
			binary.BigEndian.PutUint64(num, 0)
			slice = append(slice, num)
		}

		slice = append(slice, value...)
		buffer, err := json.Marshal(slice)

		err = pipe.Set(ctx, string(key), buffer, 0).Err()
		if err != nil {
			return err
		}

		pipe.Expire(ctx, "tx_pipelined_counter", time.Hour)
		return nil
	})

	if err != nil {
		return nil, 0, nil
	}

	return slice[1:], versionNumber, nil
}

func (r RedisDB) Del(key []byte) error {
	ctx := context.Background()

	err := r.Db.Del(ctx, string(key)).Err()
	if err != nil {
		return err
	}

	return nil
}

/*
func (r RedisDB) Replicate(key []byte) {
	ctx := context.Background()

	err := r.Db.Del(ctx, string(key)).Err()
	if err != nil {
		log.Fatal(err)
	}
}*/

func ConnectToRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.17.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
