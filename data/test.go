package data

import (
	"junk"
	"log"
	"other/import/here"
)

type thing func() error

const (
	x = 1
)

type Item struct {
	Name  string `data:"primary"`
	Make  string
	Model string
	Year  int
}

func NewItem() *Item {
	return nil
}

func (this *Item) Action() error {
	return nil
}
