# Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
# SPDX-License-Identifier: Apache-2.0

PROJDIR := $(shell readlink -f ..)
TOP_DIR := .
CUR_DIR := $(shell pwd)
PREFIX :=
DESTDIR ?= $(PREFIX)/bin

DEBUG ?=

TARGET_DIR := target
BIN_NAME := ccnp_server

CARGO := cargo

ifdef DEBUG
    release :=
    TARGET_DIR := $(TARGET_DIR)/debug
else
    release := --release
    TARGET_DIR := $(TARGET_DIR)/release
endif

TARGET := $(TARGET_DIR)/$(BIN_NAME)

test:
	$(CARGO) test

image:
	docker build \
		--build-arg http_proxy \
		--build-arg https_proxy \
		--build-arg no_proxy \
		-t ccnp-server .

build:
	$(CARGO) build $(release)

install:
	install -D -m0755 $(TARGET) $(DESTDIR)

uninstall:
	rm -f $(DESTDIR)/$(BIN_NAME)

clean:
	$(CARGO) clean