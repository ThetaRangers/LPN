package database

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisDB struct {
	Db *redis.Client
}

func (r RedisDB) Get(key []byte) [][]byte {
	ctx := context.Background()
	var slice [][]byte

	val, err := r.Db.Get(ctx, string(key)).Bytes()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(val, &slice)

	return slice
}

func (r RedisDB) Put(key []byte, value []byte) {
	ctx := context.Background()

	entry := make([][]byte, 0)
	entry = append(entry, value)

	buffer, err := json.Marshal(entry)

	err = r.Db.Set(ctx, string(key), buffer, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (r RedisDB) Append(key, value []byte) {
	ctx := context.Background()
	var slice [][]byte

	_, err := r.Db.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		val, err := pipe.Get(ctx, string(key)).Bytes()
		if err != nil {
			return nil
		}

		err = json.Unmarshal(val, &slice)
		if err != nil {
			return nil
		}

		slice = append(slice, value)
		buffer, err := json.Marshal(slice)

		err = pipe.Set(ctx, string(key), buffer, 0).Err()
		if err != nil {
			return err
		}

		pipe.Expire(ctx, "tx_pipelined_counter", time.Hour)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (r RedisDB) Del(key []byte) {
	ctx := context.Background()

	err := r.Db.Del(ctx, string(key)).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectToRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.17.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
