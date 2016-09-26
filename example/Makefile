.PHONY: deps stores

deps:
	go install github.com/zpatrick/go-sdata

stores:
	go-sdata models/user.go User --package stores > stores/user_store.go
