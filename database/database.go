package database

type Database interface {
	Get(key []byte) [][]byte
	Put(key, value []byte)
	Append(key, value []byte)
	Del(key []byte)
}
