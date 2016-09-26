# Go Simple Data
The `go-sdata` package uses code generation to quickly implement a data store for your go applications.

# Installation
```
go install github.com/zpatrick/go-sdata
```
# Motivation
When building a proof-of-concept, a small tool, or something for fun - I often find myself getting bogged down when I need to write a data persistence layer. 
In these instances, I care more about how quickly I can get my application up and running than I do building  some high-performance data store. 

Go-sdata provides you with a simple, readable, data store: 
```
package main

import (
        "log"
        "github.com/zpatrick/go-sdata/example/stores"
        "github.com/zpatrick/go-sdata/container"
)

type User struct {
	ID string `data:"primary_key"`
	Name string
}

func main() {
        fileContainer := container.NewStringFileContainer("users.json", nil)
        userStore := stores.NewUserStore(fileContainer)

        if err := userStore.Init(); err != nil {
                log.Fatal(err)
        }

        users, err := userStore.SelectAll().Execute()
        if err != nil {
                log.Fatal(err)
        }

        log.Println(users)
}
```

**Be Warned!** Go-sdata is **not** meant for high performance applications. 
As stated above, it is best suited for proof-of-concepts and other small applications. 
Please see the [Caveats](#caveats) section for more details.

# Usage
Go-sdata has 2 parts:
* A command line tool that generates a data store
* A `containers` package that is used to store/retrieve data

#### Command Line 
```
Usage: go-sdata PATH STRUCT [--package] [--template]

Arguments:
  PATH      Path to the source file
  STRUCT    Name of the source struct

Options:
  -p, --package=""    Package name for the destination file
  -t, --template=""   Path to the template file
```

See the example [Makefile](https://github.com/zpatrick/go-sdata/blob/master/example/Makefile) for reference. 

#### Containers
Containers are used to store and retrieve your data. 
They implement the [Container](https://godoc.org/github.com/zpatrick/go-sdata/container#Container) interface and can be swapped out as necessary.
The current containers are:
* [StringFileContainer](https://godoc.org/github.com/zpatrick/go-sdata/container#NewStringFileContainer) - Stores data in a human-readable JSON file 
* [ByteFileContainer](https://godoc.org/github.com/zpatrick/go-sdata/container#NewByteFileContainer) - Stores data in a byte-formatted JSON file
* [MemoryContainer](https://godoc.org/github.com/zpatrick/go-sdata/container#NewMemoryContainer) - Stores data in-memory during program execution (useful for testing)

#### Functions
Concrete examples of the following functions can be seen in the [Example](https://github.com/zpatrick/go-sdata/tree/master/example) application. 
Since the generated data store is based off of some unknown type, the `<type>` in these examples is used as a placeholder.

**Init** - Initializes the data store and container layer.
```
if err := store.Init(); err != nil {
    log.Fatal(err)
}
```

**Insert** - Inserts a new` <type>` object. The `<type>` object's primary key must be unique or else an error will be thrown. 
```
if err := store.Insert(<type>).Execute(); err != nil {
    log.Fatal(err)
}
```

**Delete** - Deletes the object with the specified primary key. Returns a boolean if the object existed.
```
existed, err := store.Delete("primary_key").Execute()
if err != nil {
    log.Fatal(err)
}

if !existed {
    log.Fatal("Item does not exist")
}
```

**SelectAll** - Returns all objects in the store
```
objects, err := store.SelectAll().Execute()
if err != nil {
    log.Fatal(err)
}
```

**SelectAllWhere** - Runs a filter on the objects before returning. 
This must be chained to a `SelectAll` query. 
The filter should return `true` if the object should be included in the results. 
```
objects, err := store.SelectAll().Where(func(t <type>) bool {return true}).Execute()
if err != nil {
    log.Fatal(err)
}
```

**FirstOrNil** - Returns the first object returned, or nil if an empty slice was returned. 
This must be chained to a `SelectAll` query. 
```
object, err := store.SelectAll().Where(func(t <type>) bool {return true}).FirstOrNil().Execute()
if err != nil {
    log.Fatal(err)
}
```

# Caveats
Go-sdata works under the following constraints:
* The struct being stored has a string primary key field (denoted by adding the field tag `"data:primary_key"`)
* The struct being stored can be marshalled/unmarshalled with the [JSON](https://golang.org/pkg/encoding/json) package
* The container used to store/retrieve data implements the following operations:
  * SelectAll()
  * Delete(primary_key)
  * Insert(primary_key, data)

Note that the container is only required to implement `SelectAll()`. 
As a result, any select operation is extremely expensive. 
This is why I only recommend using go-sdata for proof-of-concept and other small applications.
