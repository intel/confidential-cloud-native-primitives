
# Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
# SPDX-License-Identifier: Apache-2.0

PROJDIR := $(shell readlink -f ..)
TOP_DIR := .
CUR_DIR := $(shell pwd)
PREFIX :=
DESTDIR ?= $(PREFIX)/bin

DEBUG ?=

TARGET_DIR := target
BIN_NAME := quote_server

CARGO := /root/.cargo/bin/cargo

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

build:
	$(CARGO) build $(release)

install:
	install -D -m0755 $(TARGET) $(DESTDIR)

uninstall:
	rm -f $(DESTDIR)/$(BIN_NAME)

clean:
	$(CARGO) clean