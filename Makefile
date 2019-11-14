all: commit
commit:
	 go test ./...
	 golangci-lint run
	 golint ./...
	 go mod tidy
	 go install && h7t docs
	 go build -o ./plugins/csv/transformer ./plugins/csv/main.go
tools: download
	cat tools.go | grep _ | awk -F'"' '{print $$2'} | xargs -tI % go install %
download:
	 go mod download
build:
	go install
	go build -o ./plugins/csv/transformer ./plugins/csv/main.go
