.PHONY : help-default install-default build-default test-default coverage-default fmt-default mod-default archive-default lint-default generate-default tools-default runner-default snapshot-default licenses-default install-default security-default outdated-default lines-default

DEFAULTS_VERSION := 2.9.2
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(notdir $(patsubst %/,%,$(dir $(MAKEFILE_PATH))))
SRC := $(shell find . -type f -name "*.go" -not -path "./vendor/*") # used when tools are not vendor aware
MODULE := $(shell go list)
LD_VERSION := $(shell git describe --tags --abbrev=0 --dirty=-next)
LD_COMMIT := $(shell git rev-parse HEAD)
LD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_FLAGS := -s -w -X $(MODULE)/cmd.version=$(LD_VERSION) -X $(MODULE)/cmd.commit=$(LD_COMMIT) -X $(MODULE)/cmd.date=$(LD_DATE)

all: mod generate fmt test lint install

install-default: ## install the binary
	go install -ldflags="$(LD_FLAGS)" ./...

build-default: ## build the binary, ignoring vendored code if it exists
	go build -ldflags="$(LD_FLAGS)" ./...
        
test-default: ## run test with coverage
	go test -v -cover ./... -coverprofile cover.out

coverage-default: ## report on test coverage
coverage-default: test 
	goverreport -coverprofile=cover.out -sort=block -order=desc -threshold=85

fmt-default: ## organise import and format the code
	@echo "goimports && gofmt" # dont show real command as SRC could be huge
	@find $(SRC)  -type f  -name "*.go" -not -path "*/mocks/*"  -exec goimports  -w  {} \; && gofmt -s -w $(SRC)

mod-default: ## makes sure go.mod matches the source code in the module
	go mod tidy

archive-default: ## archive the third party dependencies, typically prior to generating a tagged release
	go mod vendor

lint-default: ## run golangci-lint using the configuration in .golangci.yml
	golangci-lint run

generate-default: ## go generate code
	go generate ./...
	
tools-default: ## install the project specific tools
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

runner-default: ## execute the gitlab runner using the configuration in .gitlab-ci.yml
	gitlab-runner exec docker --cache-dir /cache --docker-volumes 'cache:/cache' test

snapshot-default: ## generate a snapshot release using goreleaser
	goreleaser --snapshot --rm-dist

licenses-default: ## print list of licenses for third party software used in binary, if using repeatedly, use GITHUB_TOKEN
licenses-default: install
	golicense .approved-licenses.json $(GOPATH)/bin/$(CURRENT_DIR)

security-default: ## run go security check
security-default:
	gosec -conf .gosec.json ./...

outdated-default: ## check for outdated direct dependencies
outdated-default:
	go list -u -m -json all | go-mod-outdated -direct

lines-default: ## shorten lines longer than 100 chars, ignore generated
lines-default:
	golines --ignore-generated -m 100 -w .

help-default:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m ignore suffix -default e.g. make install \n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

%:  %-default
	@  true
  
##
## Boilerplate code, should be commmented out/deleted after initial project setup
## 
REGISTRY ?= github.com
ORIGIN ?= github.com
GROUP ?= damianoneill
PROJECT ?= h7t

boilerplate: go.mod tools.go .gitignore .golangci.yml cmd cmd/version.go cmd/completion.go pkg/domain .goreleaser.yml Dockerfile .gitlab-ci.yml .approved-licenses.json .gosec.json .git
	
.PHONY : $(PHONY) .idea

go.mod:
	@go mod init $(ORIGIN)/$(GROUP)/$(PROJECT) > /dev/null 2>&1

define TOOLS_BODY
// +build tools

package main

import (
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/mcubik/goverreport"
	_ "github.com/mitchellh/golicense"
	_ "github.com/psampaz/go-mod-outdated"
	_ "github.com/securego/gosec/cmd/gosec"
	_ "github.com/segmentio/golines"
	_ "github.com/spf13/cobra/cobra"
	_ "golang.org/x/tools/cmd/goimports"
)
endef
export TOOLS_BODY
tools.go:
	@echo "$$TOOLS_BODY" > tools.go

cmd:
	@cobra init --pkg-name $(ORIGIN)/$(GROUP)/$(PROJECT) > /dev/null 2>&1
	@rm LICENSE

pkg/domain:
	@mkdir -p pkg/{infrastructure,adapter,usecase,domain}
	@find . -type d -empty -not -path "./.git/*" -exec touch {}/.gitkeep \;

.gitignore:
	@wget -q -O .gitignore https://raw.githubusercontent.com/github/gitignore/master/Go.gitignore 
	@if ! grep -Fxq ".idea" .gitignore; then echo "\n# Exclude the Goland configuration directory\n.idea" >> .gitignore; fi 
	@if ! grep -Fxq "cover.out" .gitignore; then echo "\n# Exclude the coverage generated file \ncover.out" >> .gitignore; fi

.golangci.yml:
	@wget -q https://raw.githubusercontent.com/golangci/golangci-lint/master/.golangci.yml

define GORELEASER_BODY
env:
  - CGO_ENABLED=0
before:
  hooks:
    - make
    # - make archive
gitlab_urls:
  api: https://$(ORIGIN)/api/v4/
  download: https://$(ORIGIN)
builds:
  - id: "$(PROJECT)"
    main: ./main.go
    binary: $(PROJECT)
    ldflags:
     - -s -w -X $(ORIGIN)/$(GROUP)/$(PROJECT)/cmd.version={{.Version}} -X $(ORIGIN)/$(GROUP)/$(PROJECT)/cmd.commit={{.Commit}} -X $(ORIGIN)/$(GROUP)/$(PROJECT)/cmd.date={{.Date}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
dockers:
 - binaries:
     - $(PROJECT)
   image_templates:
     - "$(REGISTRY)/$(GROUP)/$(PROJECT):latest"
     - "$(REGISTRY)/$(GROUP)/$(PROJECT):{{ .Tag }}"
   dockerfile: Dockerfile
archives:
  - replacements:
      amd64: x86_64
release:
  gitlab:
    owner: $(GROUP)
    name: $(PROJECT)
checksum:
  name_template: "checksums.txt"
snapshot:
  name_template: "{{ .Tag }}-next"
changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
endef
export GORELEASER_BODY
.goreleaser.yml:
	@echo "$$GORELEASER_BODY" > .goreleaser.yml 
	@if ! grep -Fxq "dist" .gitignore; then echo "\n# Exclude the goreleaser artefacts directory\ndist" >> .gitignore; fi

define DOCKERFILE_BODY
FROM scratch
COPY $(PROJECT) ./
CMD ["./$(PROJECT)"]
endef
export DOCKERFILE_BODY
Dockerfile:
	@echo "$$DOCKERFILE_BODY" > Dockerfile 


define GITLABCI_BODY
.go-cache:
  variables:
    GOPATH: $$CI_PROJECT_DIR/.go
  before_script:
    - mkdir -p .go
  cache:
    paths:
      - .go/pkg/mod/

test:
  image: golang:1.14
  extends: .go-cache
  script:
    - make test
  tags:
    - docker
endef
export GITLABCI_BODY
.gitlab-ci.yml:
	@echo "$$GITLABCI_BODY" > .gitlab-ci.yml 

define VERSION_BODY
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the current build information",
	Long:  "Version, Commit and Date will be output from the Build Info.",
	Run: func(cmd *cobra.Command, args []string) {
		Version(os.Stdout)
	},
}

// Version outputs a formatted version message to the passed writer.
func Version(w io.Writer) {
	fmt.Fprintf(w, "%v, commit %v, built at %v \\n", version, commit, date)
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(versionCmd)
}
endef
export VERSION_BODY
cmd/version.go:
	@echo "$$VERSION_BODY" > cmd/version.go 


define COMPLETION_BODY
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completionTarget string

// completionCmd represents the completion command.
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion script for " + rootCmd.Use,
	Long: `Generates a shell completion script for ` + rootCmd.Use + `.
	NOTE: The current version supports Bash only.
		  This should work for *nix systems with Bash installed.

	By default, the file is written directly to /etc/bash_completion.d
	for convenience, and the command may need superuser rights, e.g.:

		sudo ` + rootCmd.Use + ` completion
	
	Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
	file-path and name.

	For e.g. on OSX with bash completion installed with brew you should 

	` + rootCmd.Use + ` completion --completionfile /etc/bash_completion.d/` + rootCmd.Use + `.sh

	Logout and in again to reload the completion scripts,
	or just source them directly:

		. /etc/bash_completion
		
	or using if using brew
	
		. $$(brew --prefix)/etc/bash_completion`,

	Run: Completion,
}

// Completion is a helper function to allow passing arguments to
// other functions (so that they can be unit tested).
func Completion(cmd *cobra.Command, args []string) {
	err := cmd.Root().GenBashCompletionFile(completionTarget)
	completion(err)
}

func completion(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Bash completion file for "+rootCmd.Use+" saved to", completionTarget)
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(completionCmd)

	completionCmd.PersistentFlags().StringVarP(&completionTarget,
		"completionfile",
		"",
		"/etc/bash_completion.d/"+rootCmd.Use+".sh",
		"completion file")
	// Required for bash-completion
	_ = completionCmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})
}
endef
export COMPLETION_BODY
cmd/completion.go:
	@echo "$$COMPLETION_BODY" > cmd/completion.go 
	

.git:
	@git init && git add . && git commit -m 'initial commit using boilerplate v$(DEFAULTS_VERSION)' .

define APPROVED_LICENSES_BODY
{
  "allow": ["MIT", "Apache-2.0", "BSD-2-CLAUSE","BSD-3-CLAUSE","MPL-2.0"]
}
endef
export APPROVED_LICENSES_BODY
.approved-licenses.json:
	@echo "$$APPROVED_LICENSES_BODY" > .approved-licenses.json 

define GOSEC_BODY
{
  "global": {
    "nosec": "enabled",
    "audit": "enabled"
  }
}
endef
export GOSEC_BODY
.gosec.json:
	@echo "$$GOSEC_BODY" > .gosec.json 
# end of boilerplate