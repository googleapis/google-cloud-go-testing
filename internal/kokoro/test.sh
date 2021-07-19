#!/bin/bash

# TODO(deklerk) Add integration tests when it's secure to do so. b/64723143

# Fail on any error
set -eo pipefail

# Display commands being run
set -x

# cd to project dir on Kokoro instance
cd github/google-cloud-go-testing

go version

# Set $GOPATH
export GOPATH="$HOME/go"
export GCGT_HOME=$GOPATH/src/github.com/googleapis/google-cloud-go-testing
export PATH="$GOPATH/bin:$PATH"
export GO111MODULE=on
mkdir -p $GCGT_HOME

# Move code into $GOPATH and get dependencies
git clone . $GCGT_HOME
cd $GCGT_HOME

try3() { eval "$*" || eval "$*" || eval "$*"; }

# All packages, including +build tools, are fetched.
try3 go mod download
./internal/kokoro/vet.sh

# Run tests and tee output to log file, to be pushed to GCS as artifact.
go test -race -v ./... 2>&1 | tee $KOKORO_ARTIFACTS_DIR/$KOKORO_GERRIT_CHANGE_NUMBER.txt
