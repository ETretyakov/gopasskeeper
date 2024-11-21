PROJECTNAME=$(shell basename "$(PWD)")

VERSION=v1.0.1
GOBASE=$(shell pwd)
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

go-build-client:
	@echo "  >  Building client"
		go build -o builds/client -ldflags "-X main.Version=$(VERSION)" cmd/client/main.go

go-build-server:
	@echo "  >  Building server"
		go build -o builds/server -ldflags "-X main.Version=$(VERSION)" cmd/server/main.go

go-build-migrate:
	@echo "  >  Building migrate"
		go build -o builds/migrate -ldflags "-X main.Version=$(VERSION)" cmd/migrate/main.go

go-build-cert:
	@echo "  >  Building cert"
		go build -o builds/cert -ldflags "-X main.Version=$(VERSION)" cmd/cert/main.go

go-build:
	bash -c "$(MAKE) go-build-client"
	bash -c "$(MAKE) go-build-server"
	bash -c "$(MAKE) go-build-migrate"
	bash -c "$(MAKE) go-build-cert"

go-client:
	@echo "  >  Satrting client"
	bash -c "$(MAKE) go-build-client"
	bash -c "builds/client"

go-server:
	@echo "  >  Satrting server"
	bash -c "$(MAKE) go-build-server"
	bash -c "builds/server"

go-migrate:
	@echo "  >  Satrting server"
	bash -c "$(MAKE) go-build-migrate"
	bash -c "builds/migrate"
