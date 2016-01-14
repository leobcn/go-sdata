# Go Simple Data
The `go-sdata` package is used to quickly implement data persistance for your go applications.

# Motivation

I often find myself spending way too much time implementing some crude form of data persistance in small go applications. The keyword here is _small_: when I'm building a proof of concept, a small tool, or something for fun - I care about how quickly I can get my code working a hell of a lot more than the performance of storing/retriveing data (especially when I know the data sets are never going to be very large).

> "Ok, I need some storage. I'll just use the `json` package to read/save from a file at first. Spend some time looking up how to do everything (load/save a file? Json marshal/unmarshal? Pretty print json?). Spend a while getting your json store to work, then add more data that need persistance. Either spend some time making your json store more robust/generic, or just copy paste a bunch of code around. Regardless, get that working. Then, down the line, you have to use a database for data persistance. OK, let me google popular golang mysql packages. Spend way too long on the readme trying to figure out the naunces of this package. Build a small little app to test their package. Can't figure it out in 5 minutes? Lookup another package. Wash, rinse, repeat until you've finally found a package that works. Crap, now I've got to find a way to make this generic enough to use my existing json stuff, or just replace all of it. Testing? Is there a `httptest` equivalent to the `sql` package? No? Crap, I knew I should have made my store use an interface...

Anyways, point is this package gives you the following:

##### It's easy to get code running
I've tried to make this tool as simple as possible to use. This package uses code generation. I've tried to lock down on struct requiresments as much as possible: just tell the tool which field to use as the primary key by giving it a tag: `data:primary_key`. 

##### It's readable
We all understand code - including compilers. Don't have to debug select statements, keep logic in the code. 

```
type User struct{
    Name string
    Age int
}
```

Get users who are 27 using sql package 
(this based off of the official documentation [example](https://golang.org/pkg/database/sql/#example_DB_Query)):
```
    age := 27
    rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
    if err != nil {
            log.Fatal(err)
    }
    defer rows.Close()
    
    users := []User{}
    for rows.Next() {
            var name string
            if err := rows.Scan(&name); err != nil {
                    log.Fatal(err)
            }
            users = append(users, User{Name: name, Age: age})
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }
    
    return users
```

Same query using `go-sdata`:
```
users, err := userStore.Select().Where(func(u User) bool { return u.Age == 27 }).Execute()
```

##### It's testable
Swap out your store with the in-memory one. 

