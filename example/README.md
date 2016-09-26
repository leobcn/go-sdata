# Go Simple Data Example
This example shows the basic functionality of a generated data store

### Generate the UserDataStore
The `Makefile` located in this directory can be used to generate the UserDataStore. 


If you don't have `go-sdata` installed, run the following to install it:
```
make deps
```

To tell go-sdata to generate the UserDataStore, run:
```
make stores
```

### Run the Application
To run the application, run:
```
go run main.go
```

You should see the following output:
```
$ go run main.go
Inserting: John Jane Aaron Susan
Deleting: Aaron
All Users: John Jane Susan
Adults: John Jane
Jane is 28 years old
```

The contents of `users.json` will reflect the current state of the data.
