package store

type Store interface {
	Init(string) error
	Select(string) (map[string][]byte, error)
	Insert(string, string, []byte) error
	Delete(string, string) error
}
