// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/spf13/cobra/cobra"
	_ "golang.org/x/lint/golint"
	// https://github.com/src-d/proteus
	// _ "github.com/golang/protobuf/protoc-gen-go"
	// _ "github.com/gogo/protobuf/..."
	// https://github.com/anjmao/go2proto
	// _ "github.com/anjmao/go2proto"
	// gRPC testing
	// _ "github.com/fullstorydev/grpcui/cmd/grpcui"
)
