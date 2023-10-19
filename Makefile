.PHONY: build clean deploy
SHELL:=/bin/bash

.ONESHELL:
build: clean
	X=$$(find ./handlers/* -type d)
	for dir in $$X ; do \
		env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$$dir $$dir/main.go ; \
	done

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose --region us-west-2

down: 
	sls remove --verbose --region us-west-2

fmt:
	gofumpt -l -w .

test:
	go test ./...

testv:
	go test ./... -v

gen:
	go generate ./...
