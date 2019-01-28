# ########################################################## #
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Check for required command tools to build or stop immediately
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=punch
VERSION=0.0.1
ARCHITECTURES=386 amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION}"

default: build

all: clean windows linux macos

windows:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=windows; export GOARCH=$(GOARCH); go build ${LDFLAGS} -o $(BINARY)-windows-$(GOARCH).exe))

linux:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=linux; export GOARCH=$(GOARCH); go build ${LDFLAGS} -o $(BINARY)-linux-$(GOARCH)))

macos:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=darwin; export GOARCH=$(GOARCH); go build ${LDFLAGS} -o $(BINARY)-darwin-$(GOARCH)))

build:
	go build ${LDFLAGS} -o ${BINARY}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete
