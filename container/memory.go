package container

import (
	"fmt"
)

type MemoryContainer struct {
	Tables map[string]map[string][]byte
}

func NewMemoryContainer() *MemoryContainer {
	return &MemoryContainer{
		Tables: map[string]map[string][]byte{},
	}
}

func (this *MemoryContainer) Init(tableID string) error {
	if _, exists := this.Tables[tableID]; !exists {
		this.Tables[tableID] = map[string][]byte{}
	}

	return nil
}

func (this *MemoryContainer) Select(tableID, query string) (map[string][]byte, error) {
	return nil, fmt.Errorf("Select is not implemented for MemoryContainer!")
}

func (this *MemoryContainer) SelectAll(tableID string) (map[string][]byte, error) {
	return this.Tables[tableID], nil
}

func (this *MemoryContainer) Insert(tableID, key string, entry []byte) error {
	if _, exists := this.Tables[tableID][key]; exists {
		return fmt.Errorf("Key '%s' already exists for table '%s'", key, tableID)
	}

	this.Tables[tableID][key] = entry

	return nil
}

func (this *MemoryContainer) Delete(tableID, key string) error {
	if _, exists := this.Tables[tableID][key]; exists {
		delete(this.Tables[tableID], key)
	}

	return nil
}
