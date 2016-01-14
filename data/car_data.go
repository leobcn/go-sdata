package data

import (
    "fmt"
    "github.com/zpatrick/go-simple-data/store"
)

type CarData struct {
    store store.Store
}

func NewCarData(store store.Store) *CarData {
    return &CarData{
        store: store,
    }
}

func (this *CarData) Init() error {
    return this.store.Init("Car")
}

type CarDataCreate struct {
    store store.Store
    data  Car
}

func (this *CarData) Create(data Car) *CarDataCreate {
    return &CarDataCreate{
        store: this.store,
        data:  data,
    }
}

func (this *CarDataCreate) Execute() error {
    return this.store.Insert("Car", this.data.Name, this.data)
}

type CarDataSelect struct {
    store  store.Store
    filter CarFilter
}

func (this *CarData) Select() *CarDataSelect {
    return &CarDataSelect{
        store: this.store,
    }
}

type CarFilter func(Car) bool

func (this *CarDataSelect) Where(filter CarFilter) *CarDataSelect {
    this.filter = filter
    return this
}

func (this *CarDataSelect) Execute() ([]Car, error) {
    results := []Car{}

    data, err := this.store.Select("Car")
    if err != nil {
        return nil, err
    }

    for _, d := range data {
        if v, ok := d.(Car); ok {
            results = append(results, v)
        } else {
            return nil, fmt.Errorf("Failed to convert '%v' to type 'Car'", v)
        }
    }

    if this.filter != nil {
        filtered := []Car{}

        for _, result := range results {
            if this.filter(result) {
                filtered = append(filtered, result)
            }
        }

        results = filtered
    }

    return results, nil
}

type CarDataDelete struct {
    store store.Store
    data  Car
}

func (this *CarData) Delete(data Car) *CarDataDelete {
    return &CarDataDelete{
        store: this.store,
        data:  data,
    }
}

func (this *CarDataDelete) Execute() error {
    return this.store.Delete("Car", this.data.Name)
}
