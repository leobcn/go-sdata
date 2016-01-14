package store

import (
	"fmt"
)

type MemoryStore struct {
	Tables map[string]map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Tables: map[string]map[string][]byte{},
	}
}

func (this *MemoryStore) Init(tableID string) error {
	if _, exists := this.Tables[tableID]; !exists {
		this.Tables[tableID] = map[string][]byte{}
	}

	return nil
}

func (this *MemoryStore) Select(tableID string) (map[string][]byte, error) {
	return this.Tables[tableID], nil
}

func (this *MemoryStore) Insert(tableID, key string, entry []byte) error {
	if _, exists := this.Tables[tableID][key]; exists {
		return fmt.Errorf("Key '%s' already exists for table '%s'", key, tableID)
	}

	this.Tables[tableID][key] = entry

	return nil
}

func (this *MemoryStore) Delete(tableID, key string) error {
	if _, exists := this.Tables[tableID][key]; exists {
		delete(this.Tables[tableID], key)
	}

	return nil
}
