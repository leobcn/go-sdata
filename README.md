# Go Simple Data
The `go-sdata` package uses code generation to quickly implement a generic data persistance layer for your go applications.

# Installation
```
go install github.com/zpatrick/go-sdata
```
# Motivation
When building a proof-of-concept, a small tool, or something for fun - I often find myself getting bogged down when I need to write a data persistance layer. 
In these instances, I care more about how quickly I can get my application up and running than I do building  some high-performance data layer. 

Go-sdata provides you with a simple, readable, data logic layer: 
```
package main

import (
        "log"
        "github.com/zpatrick/go-sdata/example/stores"
        "github.com/zpatrick/go-sdata/container"
)

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

# Usage
Go-sdata has 2 parts:
* A command line tool that generates a data persistance layer
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

See the example [Makefile](#) for reference. 

#### Containers
Containers are used to store and retrieve your data. 
They implement the [Container](#) interface and can be swapped out as necessary.
The current containers are:
* [StringFileContainer](#) - stores data in a file in human-readable json format 
* [ByteFileContainer](#) - stores data in a file in byte format (has better performance than StringFileContainer)
* [MemoryContainer](#) - stores data in-memory during program execution (useful for testing)

#### Examples
All of these examples can be seen in the [Example](#) application. 
The `<type>` references in these examples refer to the type that the data store was built from.
Stores generated The generated data layer has the following functions:

**Init** - Initializes the data logic and container layer.
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
