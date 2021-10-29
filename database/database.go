package database

type Database interface {
	Get(key []byte) ([][]byte, uint64, error)
	Put(key []byte, value [][]byte, versions ...uint64) (uint64, error)
	Append(key, value []byte) ([][]byte, uint64, error)
	Del(key []byte) error
	//Replicate(key, value [][]byte, version uint64)
}
