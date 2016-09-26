package main

import (
	"fmt"
	"github.com/zpatrick/go-sdata/container"
	"github.com/zpatrick/go-sdata/example/models"
	"github.com/zpatrick/go-sdata/example/stores"
	"log"
)

var users = []*models.User{
                {
                        ID:   "1",
                        Name: "John",
                        Age:  40,
                },
                {
                        ID:   "2",
                        Name: "Jane",
                        Age:  28,
                },
                {
                        ID:   "3",
                        Name: "Aaron",
                        Age:  12,
                },
                {
                        ID:   "4",
                        Name: "Susan",
                        Age:  9,
                },
}

func main() {
	fileContainer := container.NewStringFileContainer("users.json", nil)
	userStore := stores.NewUserStore(fileContainer)

	if err := userStore.Init(); err != nil {
		log.Fatal(err)
	}

	insert(userStore, users)
	delete(userStore, users[2])
	selectAll(userStore)
	selectAllWhere(userStore)
	firstOrNil(userStore)
}

func printNames(message string, users ...*models.User) {
	for _, user := range users {
		message += fmt.Sprintf("%s ", user.Name)
	}

	fmt.Println(message)
}

func insert(store *stores.UserStore, users []*models.User) {
	printNames("Inserting: ", users...)
	for _, user := range users {
		// errors ignored for idempotence
		store.Insert(user).Execute()
	}
}

func delete(store *stores.UserStore, user *models.User) {
	printNames("Deleting: ", user)
	existed, err := store.Delete(user.ID).Execute()
	if err != nil {
		log.Fatal(err)
	}

	if !existed {
		log.Fatal("User did not exist")
	}
}

func selectAll(store *stores.UserStore) {
	all, err := store.SelectAll().Execute()
	if err != nil {
		log.Fatal(err)
	}

	printNames("All Users: ", all...)
}

func selectAllWhere(store *stores.UserStore) {
	query := store.SelectAll().Where(func(u *models.User) bool {
		return u.Age >= 18
	})

	adults, err := query.Execute()
	if err != nil {
		log.Fatal(err)
	}

	printNames("Adults: ", adults...)
}

func firstOrNil(store *stores.UserStore) {
	query := store.SelectAll().Where(func(u *models.User) bool {
                return u.Name == "Jane"
        }).FirstOrNil()

	jane, err := query.Execute()
	if err != nil {
		log.Fatal(err)
	}

	if jane == nil {
		log.Fatal("Could not find anyone named 'Jane'")
	}

	fmt.Printf("Jane is %d years old\n", jane.Age)
}
