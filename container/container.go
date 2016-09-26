package container

// Container is an interface that can be used to store and retrieve data
type Container interface {
	// Init is used to initialize the specified table 
	Init(table string) error
	// SelectAll returns all entries in the specified table in a map.
	// The keys in the map are the primary keys, the values are the data
	SelectAll(table string) (map[string][]byte, error)
	// Insert adds the specified data to the specified table and primary key
	// If the key already exists, an error will be returned
	Insert(table, key string, data []byte) error
	// Delete will remove the entry at the specified table and key. 
	// It will return a boolean specifying if the key existed. 
	Delete(table, key string) (bool, error)
}
