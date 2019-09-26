#!/bin/bash

# Display commands being run
set -x

# Only run the linter on go1.12, since it needs type aliases (and we only care
# about its output once).
if [[ `go version` != *"go1.12"* ]]; then
    exit 0
fi

go install \
  golang.org/x/lint/golint \
  golang.org/x/tools/cmd/goimports \
  honnef.co/go/tools/cmd/staticcheck

# Fail if a dependency was added without the necessary go.mod/go.sum change
# being part of the commit.
go mod tidy
git diff go.mod | tee /dev/stderr | (! read)
git diff go.sum | tee /dev/stderr | (! read)

# Easier to debug CI.
pwd

# Look at all .go files (ignoring .pb.go files) and make sure they have a Copyright. Fail if any don't.
find . -type f -name "*.go" ! -name "*.pb.go" -exec grep -L "\(Copyright [0-9]\{4,\}\)" {} \; 2>&1 | tee /dev/stderr | (! read)
gofmt -s -d -l . 2>&1 | tee /dev/stderr | (! read)
goimports -l . 2>&1 | tee /dev/stderr | (! read)

golint ./... 2>&1 | ( \
    grep -v "should have comment or be unexported" | \
    grep -v "doc.go:17:1: package comment should be of the form" || true) | tee /dev/stderr | (! read)

staticcheck -ignore '
*:SA1019
' ./...
