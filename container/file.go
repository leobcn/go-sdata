package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type FileContainer struct {
	path  string
	mutex *sync.Mutex
}

func NewFileContainer(path string, mutex *sync.Mutex) *FileContainer {
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	return &FileContainer{
		path:  path,
		mutex: mutex,
	}
}

func (this *FileContainer) Init(tableID string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, err := os.Stat(this.path); err != nil {
		ioutil.WriteFile(this.path, []byte("{}"), os.FileMode(0777))
	}

	tables, err := this.getTables(tableID)
	if err != nil {
		return err
	}

	return this.setTables(tables)
}

func (this *FileContainer) Select(tableID, query string) (map[string][]byte, error) {
	return nil, fmt.Errorf("Select is not implemented for FileContainer!")
}

func (this *FileContainer) SelectAll(tableID string) (map[string][]byte, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(tableID)
	if err != nil {
		return nil, err
	}

	return tables[tableID], nil
}

func (this *FileContainer) Insert(tableID, key string, entry []byte) error {
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

func (this *FileContainer) Delete(tableID, key string) (bool, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tables, err := this.getTables(tableID)
	if err != nil {
		return false, err
	}

	if _, exists := tables[tableID][key]; !exists {
		return false, nil
	}

	delete(tables[tableID], key)
	return true, this.setTables(tables)
}

func (this *FileContainer) getTables(tableID string) (map[string]map[string][]byte, error) {
	bytes, err := ioutil.ReadFile(this.path)
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

func (this *FileContainer) setTables(tables map[string]map[string][]byte) error {
	bytes, err := json.MarshalIndent(tables, "", "    ")
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(this.path, bytes, os.FileMode(0777))
}
