MODULE=gouploadserver
VERSION=v0.0.2-alpha.0
BUILDTIME=$(shell date +"%Y-%m-%dT%T%z")
# FIXME add LDFLAGS
# LDFLAGS= -ldflags '-X ...version=$(VERSION) -X ....buildTime=$(BUILDTIME)'

.PHONY: default
default: build

.PHONY: clearbin
clearbin:
	rm -rf ./bin

build: main.go clearbin
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(MODULE)-linux-amd64 ./main.go
	GOOS=windows GOARCH=amd64 go build -v -o ./bin/$(MODULE)-windows-amd64.exe ./main.go
	GOOS=darwin GOARCH=amd64 go build -v -o ./bin/$(MODULE)-darwin-amd64 ./main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o ./bin/$(MODULE)-alpine-linux-amd64 ./main.go
	GOOS=linux GOARCH=arm GOARM=7 go build -v -o ./bin/$(MODULE)-linux-armv7 ./main.go

buildzip: main.go clearbin
	rm -rf ./dist && mkdir -p ./dist
	
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(MODULE) ./main.go
	zip -r ./dist/$(MODULE)-$(VERSION)-linux-amd64.zip ./bin/$(MODULE)
	
	GOOS=windows GOARCH=amd64 go build -v -o ./bin/$(MODULE).exe ./main.go
	zip -r ./dist/$(MODULE)-$(VERSION)-windows-amd64.zip ./bin/$(MODULE).exe
	
	GOOS=darwin GOARCH=amd64 go build -v -o ./bin/$(MODULE) ./main.go
	zip -r ./dist/$(MODULE)-$(VERSION)-darwin-amd64.zip ./bin/$(MODULE)
	
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o ./bin/$(MODULE) ./main.go
	zip -r ./dist/$(MODULE)-$(VERSION)-alpine-linux-amd64.zip ./bin/$(MODULE)
	
	GOOS=linux GOARCH=arm GOARM=7 go build -v -o ./bin/$(MODULE) ./main.go
	zip -r ./dist/$(MODULE)-$(VERSION)-linux-armv7.zip ./bin/$(MODULE)

.PHONY: cross
cross: main.go clearbin
	go build -v -o ./bin/$(MODULE) ./main.go

.PHONY: run
run: cross
	./bin/$(MODULE)

install: main.go
	go list -f '{{.Target}}'
	go install
