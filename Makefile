GOCMD=go
GOINSTALL=$(GOCMD) install

all: tools
tools: download
	cat tools.go | grep _ | awk -F'"' '{print $$2'} | xargs -tI % $(GOINSTALL) %
download:
	$(GOCMD) mod download
