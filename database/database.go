package database

type Database interface {
	Get(key []byte) ([][]byte, error)
	Put(key []byte, value [][]byte) error
	Append(key []byte, value [][]byte) ([][]byte, error)
	Del(key []byte) error
}
