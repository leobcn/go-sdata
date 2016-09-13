package container

type Container interface {
	Init(table string) error
	SelectAll(table string) (map[string][]byte, error)
	Insert(table, key string, data []byte) error
	Delete(table, key string) (bool, error)
}
