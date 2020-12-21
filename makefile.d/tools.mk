# Common tools like go-compiler

## go
tools: build/tools/go build/tools/rice
build/tools/go:
	@which $(notdir $@)
build/tools/rice: build/tools/go
	@which $(notdir $@) || (go get -u github.com/GeertJohan/go.rice/rice)


## docker
build/tools/docker:
	@which $(notdir $@)


## nodejs
tools: build/tools/npm build/tools/yarn
build/tools/npm:
	@which $(notdir $@)
build/tools/yarn: build/tools/npm
	@which $(notdir $@) || (npm install -g yarn)

#tools: build/tools/rollup
#build/tools/rollup: build/tools/npm
#	@which $(notdir $@) || (npm install -g rollup && npm install -g '@rollup/plugin-node-resolve')
#
#tools: build/tools/lessc
#build/tools/lessc: build/tools/npm
#	@which $(notdir $@) || (npm install -g less)


## grpc
tools: build/tools/protoc build/tools/protoc-gen-go build/tools/protoc-gen-go-grpc
build/tools/protoc:
	@which $(notdir $@) || (echo "see https://grpc.io/docs/protoc-installation/")
build/tools/protoc-gen-go: build/tools/go
	@which $(notdir $@) || (go get -u google.golang.org/protobuf/cmd/protoc-gen-go)
build/tools/protoc-gen-go-grpc: build/tools/go
	@which $(notdir $@) || (go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc)

#tools: build/tools/protoc-gen-grpc-gateway
build/tools/protoc-gen-grpc-gateway: build/tools/go
	@which $(notdir $@) || (go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway)

#tools: build/tools/protoc-gen-openapiv2
build/tools/protoc-gen-openapiv2: build/tools/go
	@which $(notdir $@) || (go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2)



.PHONY: tools
