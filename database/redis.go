package database

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Db *redis.Client
}

func (r RedisDB) Get(key []byte) ([][]byte, error) {
	ctx := context.Background()
	var slice [][]byte

	val, err := r.Db.Get(ctx, string(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (r RedisDB) Put(key []byte, value [][]byte) error {
	ctx := context.Background()

	buffer, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Db.Set(ctx, string(key), buffer, 0).Err()
}

func (r RedisDB) Append(key []byte, value [][]byte) ([][]byte, error) {
	ctx := context.Background()
	var slice [][]byte

	txnPut := func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, string(key)).Bytes()
		if err != nil && err != redis.Nil {
			return err
		}

		if len(val) != 0 {
			err = json.Unmarshal(val, &slice)
			if err != nil {
				return err
			}
		}

		slice = append(slice, value...)
		buffer, err := json.Marshal(slice)
		if err != nil {
			return err
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, string(key), buffer, 0).Err()
		})

		return nil
	}

	for {
		err := r.Db.Watch(ctx, txnPut, string(key))
		if err == nil {
			// Success.
			return slice, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return nil, err
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

func (r RedisDB) Migrate(key []byte) ([][]byte, error) {
	ctx := context.Background()
	var slice [][]byte

	txn := func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, string(key)).Bytes()
		if err == redis.Nil {
			return nil
		} else if err != nil {
			return err
		}

		err = json.Unmarshal(val, &slice)
		if err != nil {
			return err
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Del(ctx, string(key)).Err()
		})

		return nil
	}

	for {
		err := r.Db.Watch(ctx, txn, string(key))
		if err == nil {
			// Success.
			return slice, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return nil, err
	}
}

func (r RedisDB) DeleteExcept(keys []string) error {
	ctx := context.Background()

	var keysDb []string

	txnDel := func(tx *redis.Tx) error {
		var cursor uint64
		var err error

		keysDb, _, err = tx.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}

		for _, k := range keysDb {
			if !Contains(keys, k) {
				tx.Del(ctx, k)
			}
		}

		return nil
	}

	for {
		err := r.Db.Watch(ctx, txnDel, keysDb...)
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

func (r RedisDB) GetAllKeys() ([]string, []string) {
	ctx := context.Background()
	var keysDb []string

	txnDel := func(tx *redis.Tx) error {
		var cursor uint64
		var err error

		keysDb, _, err = tx.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}

		return nil
	}

	for {
		err := r.Db.Watch(ctx, txnDel, keysDb...)
		if err == nil {
			// Success.
			return keysDb, keysDb
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return []string{}, []string{}
	}
}

func ConnectToRedis(addr string) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
