package database

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
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

func (r RedisDB) Put(key []byte, value [][]byte) (uint64, error) {
	ctx := context.Background()
	var versionNum uint64
	var slice [][]byte

	txnPut := func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, string(key)).Bytes()
		if err == redis.Nil {
			versionNum = 0
		} else if err != nil {
			return err
		} else {
			err = json.Unmarshal(val, &slice)
			if err != nil {
				return err
			}

			versionNum = binary.BigEndian.Uint64(slice[0])
			versionNum++
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

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, string(key), buffer, 0).Err()
		})

		return nil
	}

	// Let's just hope race conditions will not happen forever
	for {
		err := r.Db.Watch(ctx, txnPut, string(key))
		if err == nil {
			// Success.
			return versionNum, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return 0, err
	}
}

func (r RedisDB) Replicate(key []byte, value [][]byte, version uint64) error {
	ctx := context.Background()
	var versionNum uint64
	var slice [][]byte

	txnPut := func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, string(key)).Bytes()

		if err != nil && err != redis.Nil {
			return err
		} else if err != redis.Nil {

			err = json.Unmarshal(val, &slice)
			if err != nil {
				return err
			}

			versionNum = binary.BigEndian.Uint64(slice[0])

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

		buffer, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, string(key), buffer, 0).Err()
		})

		return nil
	}

	// Let's just hope race conditions will not happen forever
	for {
		err := r.Db.Watch(ctx, txnPut, string(key))
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}
}

func (r RedisDB) Append(key []byte, value [][]byte) ([][]byte, uint64, error) {
	ctx := context.Background()
	var slice [][]byte
	var versionNumber uint64
	var num = make([]byte, 8)

	txnPut := func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, string(key)).Bytes()
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

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, string(key), buffer, 0).Err()
		})

		return nil
	}

	for {
		err := r.Db.Watch(ctx, txnPut, string(key))
		if err == nil {
			// Success.
			return slice[1:], versionNumber, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return nil, 0, err
	}

}

func (r RedisDB) Del(key []byte) error {
	ctx := context.Background()

	err := r.Db.Del(ctx, string(key)).Err()
	if err != nil {
		return err
	}

	return nil
}

func ConnectToRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.17.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
