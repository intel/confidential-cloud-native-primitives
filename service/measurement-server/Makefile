# Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
# SPDX-License-Identifier: Apache-2.0

export GO111MODULE=on

all: clean
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64
	@go build -tags netgo -o ./measurement-server ./server/server.go

# The following is done this way as each patch on CI runs build and each merge runs deploy. So for build we don't need to build binary and hence
# no need to create a static binary with additional flags. However, for generating binary, additional build flags are necessary. This if used with
# mock plugin errors out for unit tests. So the seperation avoids the error.

build: clean test cover
deploy: build

.PHONY: test
test: clean
	# remove race condition as workaround
	# @go test -race ./...
	@go test ./...

format:
	@go fmt ./...

clean:
	@find . -name "*so" -delete
	@rm -f measurement-server coverage.html coverage.out

.PHONY: cover
cover:
	@go test -race ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
