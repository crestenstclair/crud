.PHONY: build clean deploy

build:
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose --region us-west-2

down: 
	sls remove --verbose --region us-west-2

fmt:
	gofumpt -l -w .

test:
	go test ./... -v

gen:
	go generate ./...
