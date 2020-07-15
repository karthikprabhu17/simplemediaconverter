
PROJECTNAME = $(shell basename "$(PWD)")
BASEDIR = $(shell pwd)
GOFILES = $(wildcard $(BASEDIR)/*.go)
GOBIN = $(BASEDIR)/bin
GOPATH = $(shell go env GOPATH)

pre-req:
	@go version
	if [ "$$?" != "0" ]; then \
		echo "You dont seem to have go installed.. Please install that first"; \
		/bin/false; \
	fi; \

build: pre-req
	@echo "Building ..."
	go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

install: build
	@echo "Installing in $(GOPATH)"
	@install -D -m 755 $(GOBIN)/$(PROJECTNAME) $(GOPATH)/bin/$(PROJECTNAME)

clean:
	@echo "Cleaning up build files ..."
	rm -rf $(BASEDIR)/bin

PHONY += pre-req