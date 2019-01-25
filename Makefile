# ########################################################## #
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=punch
VERSION=0.0.1
BUILD=`git rev-parse HEAD`
ARCHITECTURES=386 amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean windows linux macos

windows:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=windows; export GOARCH=$(GOARCH); go build -o $(BINARY)-windows-$(GOARCH).exe))

linux:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=linux; export GOARCH=$(GOARCH); go build -o $(BINARY)-linux-$(GOARCH)))

macos:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=darwin; export GOARCH=$(GOARCH); go build -o $(BINARY)-darwin-$(GOARCH)))

build:
	go build ${LDFLAGS} -o ${BINARY}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean install build_all all