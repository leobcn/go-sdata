package models

type User struct {
	ID   string `data:"primary_key"`
	Name string
	Age  int
}
