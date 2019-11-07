GOCMD=go
GOINSTALL=$(GOCMD) install

all: commit
commit: 
	 golangci-lint run
	 golint ./...
	 go mod tidy
tools: download
	cat tools.go | grep _ | awk -F'"' '{print $$2'} | xargs -tI % $(GOINSTALL) %
download:
	$(GOCMD) mod download
