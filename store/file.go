package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type FileStore struct {
	path  string
	mutex *sync.Mutex
}

func NewFileStore(path string, mutex *sync.Mutex) *FileStore {
	return &FileStore{
		path:  path,
		mutex: mutex,
	}
}

func (this *FileStore) Init(tableID string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, err := os.Stat(this.path); err != nil {
		ioutil.WriteFile(this.path, []byte("{}"), os.FileMode(0777))
	}

	tables, err := this.getTables(tableID)
	if err != nil {
		return err
	}

	if tables == nil {
		tables = map[string]map[string][]byte{}
	}

	if _, ok := tables[tableID]; !ok {
		tables[tableID] = map[string][]byte{}
	}

	return this.setTables(tables)
}

func (this *FileStore) Select(tableID string) (map[string][]byte, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(tableID)
	if err != nil {
		return nil, err
	}

	return tables[tableID], nil
}

func (this *FileStore) Insert(tableID, key string, entry []byte) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(tableID)
	if err != nil {
		return err
	}

	if _, exists := tables[tableID][key]; exists {
		return fmt.Errorf("Key '%s' already exists for table '%s'", key, tableID)
	}

	tables[tableID][key] = entry

	return this.setTables(tables)
}

func (this *FileStore) Delete(tableID, key string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(tableID)
	if err != nil {
		return err
	}

	if _, exists := tables[tableID][key]; exists {
		delete(tables[tableID], key)
	}

	return this.setTables(tables)
}

func (this *FileStore) getTables(tableID string) (map[string]map[string][]byte, error) {
	bytes, err := ioutil.ReadFile(this.path)
	if err != nil {
		return nil, err
	}

	var tables map[string]map[string][]byte
	if err := json.Unmarshal(bytes, &tables); err != nil {
		return nil, err
	}

	return tables, nil
}

func (this *FileStore) setTables(tables map[string]map[string][]byte) error {
	bytes, err := json.MarshalIndent(tables, "", "    ")
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(this.path, bytes, os.FileMode(0777))
}
