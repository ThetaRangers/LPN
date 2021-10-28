package database

type Database interface {
	Get(key []byte) ([][]byte, uint64)
	Put(key, value []byte)
	Append(key, value []byte)
	Del(key []byte)
	Replicate(key, value [][]byte, version uint64)
}
