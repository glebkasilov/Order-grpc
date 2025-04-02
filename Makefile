# Makefile for the project

# Define the Go compiler
GO = go

# Define the Go build command
BUILD = $(GO) build

# Define the Go test command
TEST = $(GO) test

# Define the Go run command
RUN = $(GO) run

# Define the project's main package
MAIN_PKG = cmd/main.go

# Define the project's test packages
TEST_PKGS = $(shell $(GO) list ./... | grep -v /vendor/)

# Define the project's dependencies
DEPENDENCIES = $(shell $(GO) list -f '{{ join .Deps "\n" }}' ./...)

# Define the build target
# protoc --go_out=. --go_opt=paths=./api/order \
		--go-grpc_out=. --go-grpc_opt=paths= \
		routeguide/route_guide.proto

build:
	protoc --go_out=./pkg/api/test --go_opt=paths=source_relative \
			--go-grpc_out=./pkg/api/test --go-grpc_opt=paths=source_relative \
			./api/*.proto
	$(BUILD) -o bin/main $(MAIN_PKG)

build-bin:
	$(BUILD) -o bin/main $(MAIN_PKG)

build-proto:
	protoc -I ./ -I ./google \
		--go_out=./pkg/api/test --go_opt=paths=source_relative \
		--go-grpc_out=./pkg/api/test --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./pkg/api/test --grpc-gateway_opt=paths=source_relative \
		./api/order.proto

# Define the test target
test:
	$(TEST) $(TEST_PKGS)

# Define the run target
run:
	$(RUN) $(MAIN_PKG)

# Define the clean target
clean:
	rm -f bin/main
	rmdir bin

	rm -f pkg/api/test/api/order_grpc.pb.go
	rm -f pkg/api/test/api/order.pb.go
	rm -f pkg/api/tets/api/order.pb.gw.go

clean-bin:
	rm -f bin/main
	rmdir bin

clean-proto:
	rm -f pkg/api/test/api/order_grpc.pb.go
	rm -f pkg/api/test/api/order.pb.go
	rm -f pkg/api/tets/api/order.pb.gw.go

# Define the dependencies target
deps:
	$(GO) get -u $(DEPENDENCIES)

# Define the default target
default: build

# Define the all target
all: build test