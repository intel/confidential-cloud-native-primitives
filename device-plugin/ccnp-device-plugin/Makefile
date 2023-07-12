
# SPDX-license-identifier: Apache-2.0
##############################################################################
# Copyright (c) 2023 Intel Corporation
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

export GO111MODULE=on

.PHONY: build deploy

build:
	CGO_ENABLED=0 GOOS=linux 
	@go build -a -installsuffix cgo -o build/ccnp-device-plugin cmd/server/app.go

deploy:
	helm install ccnp-device-plugin deploy/helm/ccnp-device-plugin 

clean:
	@rm -f build
