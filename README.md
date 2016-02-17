# Go Simple Data
The `go-sdata` package is used to quickly implement data persistance for your go applications.

# Motivation
I often find myself spending way too much time implementing some crude form of data persistance in small go applications. The keyword here is _small_: when I'm building a proof of concept, a small tool, or something for fun - I care about how quickly I can get my code working a hell of a lot more than the performance of storing/retriveing data (especially when I know the data sets are never going to be very large).

##### It's easy to get code running
This package uses code generation to generate 

##### It's readable, generic, data logic layer
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

