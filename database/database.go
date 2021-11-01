package database

type Database interface {
	Get(key []byte) ([][]byte, uint64, error)
	Put(key []byte, value [][]byte) (uint64, error)
	Append(key []byte, value [][]byte) ([][]byte, uint64, error)
	Del(key []byte) error
	Replicate(key []byte, value [][]byte, version uint64) error
}
