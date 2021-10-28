package database

type Database interface {
	Get(key []byte) ([][]byte, uint64)
	Put(key []byte, value [][]byte, versions ...uint64) uint64
	Append(key, value []byte) ([][]byte, uint64)
	Del(key []byte)
	//Replicate(key, value [][]byte, version uint64)
}
