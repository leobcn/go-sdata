package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// A FileContainer stores data in a JSON file
type FileContainer struct {
	path      string
	mutex     *sync.Mutex
	getTables func(path string, tableID string) (map[string]map[string][]byte, error)
	setTables func(path string, tables map[string]map[string][]byte) error
}

// A ByteFileContainer stores data in a file in byte format
// This has better performance than a StringFileContainer
// If multiple stores are using the same container, they should share the same mutex object
func NewByteFileContainer(path string, mutex *sync.Mutex) *FileContainer {
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	return &FileContainer{
		path:      path,
		mutex:     mutex,
		getTables: getByteTables,
		setTables: setByteTables,
	}
}

// A StringFileContainer stores data in human-readable format
//  If multiple stores are using the same container, they should share the same mutex object
func NewStringFileContainer(path string, mutex *sync.Mutex) *FileContainer {
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	return &FileContainer{
		path:      path,
		mutex:     mutex,
		getTables: getStringTables,
		setTables: setStringTables,
	}
}

func (this *FileContainer) Init(tableID string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	file, err := os.Stat(this.path)
	if err != nil || file.Size() == 0 {
		if err := ioutil.WriteFile(this.path, []byte("{}"), os.FileMode(0777)); err != nil {
			return err
		}
	}

	tables, err := this.getTables(this.path, tableID)
	if err != nil {
		return err
	}

	return this.setTables(this.path, tables)
}

func (this *FileContainer) SelectAll(tableID string) (map[string][]byte, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(this.path, tableID)
	if err != nil {
		return nil, err
	}

	return tables[tableID], nil
}

func (this *FileContainer) Insert(tableID, key string, entry []byte) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(this.path, tableID)
	if err != nil {
		return err
	}

	if _, exists := tables[tableID][key]; exists {
		return fmt.Errorf("Key '%s' already exists for table '%s'", key, tableID)
	}

	tables[tableID][key] = entry

	return this.setTables(this.path, tables)
}

func (this *FileContainer) Delete(tableID, key string) (bool, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(this.path, tableID)
	if err != nil {
		return false, err
	}

	if _, exists := tables[tableID][key]; !exists {
		return false, nil
	}

	delete(tables[tableID], key)
	return true, this.setTables(this.path, tables)
}

func getByteTables(path, tableID string) (map[string]map[string][]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tables map[string]map[string][]byte
	if err := json.Unmarshal(bytes, &tables); err != nil {
		return nil, err
	}

	if tables == nil {
		tables = map[string]map[string][]byte{}
	}

	if _, ok := tables[tableID]; !ok {
		tables[tableID] = map[string][]byte{}
	}

	return tables, nil
}

func setByteTables(path string, tables map[string]map[string][]byte) error {
	bytes, err := json.Marshal(tables)
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(path, bytes, os.FileMode(0777))
}

func getStringTables(path, tableID string) (map[string]map[string][]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var decodedTables map[string]map[string]interface{}
	if err := json.Unmarshal(bytes, &decodedTables); err != nil {
		return nil, err
	}

	if decodedTables == nil {
		decodedTables = map[string]map[string]interface{}{}
	}

	encodedTables := map[string]map[string][]byte{}
	for tableID, tableKeys := range decodedTables {
		for key, value := range tableKeys {
			if _, ok := encodedTables[tableID]; !ok {
				encodedTables[tableID] = map[string][]byte{}
			}

			bytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			encodedTables[tableID][key] = bytes
		}
	}

	if _, ok := encodedTables[tableID]; !ok {
		encodedTables[tableID] = map[string][]byte{}
	}

	return encodedTables, nil
}

func setStringTables(path string, encodedTables map[string]map[string][]byte) error {
	decodedTables := map[string]map[string]interface{}{}
	for tableID, tableKeys := range encodedTables {
		for key, value := range tableKeys {
			if _, ok := decodedTables[tableID]; !ok {
				decodedTables[tableID] = map[string]interface{}{}
			}

			var item interface{}
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}

			decodedTables[tableID][key] = item
		}
	}

	bytes, err := json.MarshalIndent(decodedTables, "", "    ")
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(path, bytes, os.FileMode(0777))
}
