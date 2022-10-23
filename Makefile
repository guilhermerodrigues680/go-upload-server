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
	GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o ./bin/$(MODULE)-linux-amd64 ./main.go
	GOOS=linux GOARCH=arm64 go build -v -ldflags="-s -w" -o ./bin/$(MODULE)-linux-arm64 ./main.go
	GOOS=windows GOARCH=amd64 go build -v -ldflags="-s -w" -o ./bin/$(MODULE)-windows-amd64.exe ./main.go
	GOOS=darwin GOARCH=amd64 go build -v -ldflags="-s -w" -o ./bin/$(MODULE)-darwin-amd64 ./main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags="-s -w" -o ./bin/$(MODULE)-alpine-linux-amd64 ./main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -v -ldflags="-s -w" -o ./bin/$(MODULE)-alpine-linux-arm64 ./main.go

buildzip: main.go build
	rm -rf ./dist && mkdir -p ./dist
	
	cp -fp ./bin/$(MODULE)-linux-amd64 ./bin/$(MODULE)
	zip -r ./dist/$(MODULE)-$(VERSION)-linux-amd64.zip ./bin/$(MODULE)
	rm -f ./bin/$(MODULE)

	cp -fp ./bin/$(MODULE)-linux-arm64 ./bin/$(MODULE)
	zip -r ./dist/$(MODULE)-$(VERSION)-linux-arm64.zip ./bin/$(MODULE)
	rm -f ./bin/$(MODULE)

	cp -fp ./bin/$(MODULE)-darwin-amd64 ./bin/$(MODULE)
	zip -r ./dist/$(MODULE)-$(VERSION)-darwin-amd64.zip ./bin/$(MODULE)
	rm -f ./bin/$(MODULE)

	cp -fp ./bin/$(MODULE)-windows-amd64.exe ./bin/$(MODULE).exe
	zip -r ./dist/$(MODULE)-$(VERSION)-windows-amd64.zip ./bin/$(MODULE).exe
	rm -f ./bin/$(MODULE).exe

	cp -fp ./bin/$(MODULE)-alpine-linux-amd64 ./bin/$(MODULE)
	zip -r ./dist/$(MODULE)-$(VERSION)-alpine-linux-amd64.zip ./bin/$(MODULE)
	rm -f ./bin/$(MODULE)

	cp -fp ./bin/$(MODULE)-alpine-linux-arm64 ./bin/$(MODULE)
	zip -r ./dist/$(MODULE)-$(VERSION)-alpine-linux-arm64.zip ./bin/$(MODULE)
	rm -f ./bin/$(MODULE)

.PHONY: cross
cross: main.go clearbin
	go build -v -o ./bin/$(MODULE) ./main.go

.PHONY: run
run: cross
	./bin/$(MODULE)

install: main.go
	go list -f '{{.Target}}'
	go install
