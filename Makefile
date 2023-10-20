.PHONY: build clean deploy down fmt test tstv gen lint
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
	go fmt ./...

test:
	go test ./...

testv:
	go test ./... -v

gen:
	go generate ./...

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.54.2 golangci-lint run -v
