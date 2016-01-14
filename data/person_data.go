package data

import (
    "fmt"
    "github.com/zpatrick/go-simple-data/store"
)

type PersonData struct {
    store store.Store
}

func NewPersonData(store store.Store) *PersonData {
    return &PersonData{
        store: store,
    }
}

func (this *PersonData) Init() error {
    return this.store.Init("Person")
}

type PersonDataCreate struct {
    store store.Store
    data  Person
}

func (this *PersonData) Create(data Person) *PersonDataCreate {
    return &PersonDataCreate{
        store: this.store,
        data:  data,
    }
}

func (this *PersonDataCreate) Execute() error {
    return this.store.Insert("Person", this.data.Name, this.data)
}

type PersonDataSelect struct {
    store  store.Store
    filter PersonFilter
}

func (this *PersonData) Select() *PersonDataSelect {
    return &PersonDataSelect{
        store: this.store,
    }
}

type PersonFilter func(Person) bool

func (this *PersonDataSelect) Where(filter PersonFilter) *PersonDataSelect {
    this.filter = filter
    return this
}

func (this *PersonDataSelect) Execute() ([]Person, error) {
    results := []Person{}

    data, err := this.store.Select("Person")
    if err != nil {
        return nil, err
    }

    for _, d := range data {
        if v, ok := d.(Person); ok {
            results = append(results, v)
        } else {
            return nil, fmt.Errorf("Failed to convert '%v' to type 'Person'", v)
        }
    }

    if this.filter != nil {
        filtered := []Person{}

        for _, result := range results {
            if this.filter(result) {
                filtered = append(filtered, result)
            }
        }

        results = filtered
    }

    return results, nil
}

type PersonDataDelete struct {
    store store.Store
    data  Person
}

func (this *PersonData) Delete(data Person) *PersonDataDelete {
    return &PersonDataDelete{
        store: this.store,
        data:  data,
    }
}

func (this *PersonDataDelete) Execute() error {
    return this.store.Delete("Person", this.data.Name)
}
