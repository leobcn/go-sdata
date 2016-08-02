package container

import (
	"fmt"
	"sync"
)

type MemoryContainer struct {
	Tables map[string]map[string][]byte
	mutex  *sync.Mutex
}

func NewMemoryContainer() *MemoryContainer {
	return &MemoryContainer{
		Tables: map[string]map[string][]byte{},
		mutex:  &sync.Mutex{},
	}
}

func (this *MemoryContainer) Init(tableID string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, exists := this.Tables[tableID]; !exists {
		this.Tables[tableID] = map[string][]byte{}
	}

	return nil
}

func (this *MemoryContainer) Select(tableID, query string) (map[string][]byte, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return nil, fmt.Errorf("Select is not implemented for MemoryContainer!")
}

func (this *MemoryContainer) SelectAll(tableID string) (map[string][]byte, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.Tables[tableID], nil
}

func (this *MemoryContainer) Insert(tableID, key string, entry []byte) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, exists := this.Tables[tableID][key]; exists {
		return fmt.Errorf("Key '%s' already exists for table '%s'", key, tableID)
	}

	this.Tables[tableID][key] = entry

	return nil
}

func (this *MemoryContainer) Delete(tableID, key string) (bool, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, exists := this.Tables[tableID][key]; !exists {
		return false, nil
	}

	delete(this.Tables[tableID], key)
	return true, nil
}
