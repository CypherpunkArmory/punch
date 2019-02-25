# ########################################################## #
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Check for required command tools to build or stop immediately
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=punch
VERSION=
ROLLBAR_TOKEN=
ARCHITECTURES=386 amd64

#these get replaced later by os specific build rules
GOOS=linux
GOEXT=

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X github.com/cypherpunkarmory/punch/cmd.version=$(VERSION) -X github.com/cypherpunkarmory/punch/cmd.rollbarToken=$(ROLLBAR_TOKEN) -s -w"

default: build

all: clean windows linux macos

build-os:
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build ${LDFLAGS} -o $(BINARY)-$(GOOS)-$(GOARCH)$(GOEXT); mkdir -p gh-pages/static/$(GOOS)/$(GOARCH); cp $(BINARY)-$(GOOS)-$(GOARCH)$(GOEXT) gh-pages/static/$(GOOS)/$(GOARCH)/$(BINARY)$(GOEXT)))

windows: GOOS=windows
windows: GOEXT=.exe
windows: build-os

linux: GOOS=linux
linux: GOEXT=
linux: build-os

macos: GOOS=darwin
macos: GOEXT=
macos: build-os

build:
	go build ${LDFLAGS} -o ${BINARY}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete
