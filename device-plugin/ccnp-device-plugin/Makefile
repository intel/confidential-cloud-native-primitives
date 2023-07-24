
export GO111MODULE=on

.PHONY: build deploy

build:
	CGO_ENABLED=0 GOOS=linux 
	@go build -a -installsuffix cgo -o build/ccnp-device-plugin cmd/server/app.go

deploy:
	helm install ccnp-device-plugin deploy/helm/ccnp-device-plugin 

clean:
	@rm -f build
