package stores

// Automatically generated by go-sdata. DO NOT EDIT!

import (
    "encoding/json"
    "github.com/zpatrick/go-sdata/container"
     "github.com/zpatrick/go-sdata/example/models" 
)

type UserStore struct {
    container container.Container
    table     string
}

func NewUserStore(container container.Container) *UserStore {
    return &UserStore{
        container: container,
        table:     "models.User",
    }
}

func (this *UserStore) Init() error {
    return this.container.Init(this.table)
}

type UserStoreInsert struct {
    *UserStore
    data *models.User
}

func (this *UserStore) Insert(data *models.User) *UserStoreInsert {
    return &UserStoreInsert{
        UserStore: this,
        data:       data,
    }
}

func (this *UserStoreInsert) Execute() error {
    bytes, err := json.Marshal(this.data)
    if err != nil {
        return err
    }

    return this.container.Insert(this.table, this.data.ID, bytes)
}

type UserStoreSelectAll struct {
    *UserStore
    filter UserFilter
}

func (this *UserStore) SelectAll() *UserStoreSelectAll {
    return &UserStoreSelectAll{
        UserStore: this,
    }
}

type UserFilter func(*models.User) bool

func (this *UserStoreSelectAll) Where(filter UserFilter) *UserStoreSelectAll {
    this.filter = filter
    return this
}

func (this *UserStoreSelectAll) Execute() ([]*models.User, error) {
    data, err :=  this.container.SelectAll(this.table)
    if err != nil {
        return nil, err
    }

    results := []*models.User{}
    for _, d := range data {
        var value *models.User

        if err := json.Unmarshal(d, &value); err != nil {
            return nil, err
        }

		if this.filter == nil || this.filter(value) {
			results = append(results, value)
		}
    }

    return results, nil
}

type UserStoreSelectFirst struct {
    *UserStoreSelectAll
}

func (this *UserStoreSelectAll) FirstOrNil() *UserStoreSelectFirst {
    return &UserStoreSelectFirst{
        UserStoreSelectAll: this,
    }
}

func (this *UserStoreSelectFirst) Execute() (*models.User, error) {
    results, err := this.UserStoreSelectAll.Execute()
    if err != nil {
        return nil, err
    }

    if len(results) > 0 {
        return results[0], nil
    }

    return nil, nil
}

type UserStoreDelete struct {
    *UserStore
    key string
}

func (this *UserStore) Delete(key string) *UserStoreDelete {
    return &UserStoreDelete{
        UserStore: this,
        key:        key,
    }
}

func (this *UserStoreDelete) Execute() (bool, error) {
    return this.container.Delete(this.table, this.key)
}
