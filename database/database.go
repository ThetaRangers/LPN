package database

type Database interface {
	Get(key []byte) ([][]byte, error)
	Put(key []byte, value [][]byte) error
	Append(key []byte, value [][]byte) ([][]byte, error)
	Del(key []byte) error
	DeleteExcept([]string) error
	GetAllKeys() ([]string, []string)
}

func Contains(slice []string, string string) bool {
	for _, x := range slice {
		if x == string {
			return true
		}
	}

	return false
}
