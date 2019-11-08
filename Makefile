all: commit
commit: 
	 golangci-lint run
	 golint ./...
	 go mod tidy
	 go install && h7t docs
tools: download
	cat tools.go | grep _ | awk -F'"' '{print $$2'} | xargs -tI % go install %
download:
	 go mod download
build:
	go install
	go build -o ./plugins/delimited/transformer ./plugins/delimited/main.go
