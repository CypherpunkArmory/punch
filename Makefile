# ########################################################## #
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Check for required command tools to build or stop immediately
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=punch
VERSION=
ROLLBAR_TOKEN=
ARCHITECTURES=386 amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X github.com/cypherpunkarmory/punch/cmd.version=$(VERSION) -X github.com/cypherpunkarmory/punch/cmd.rollbarToken=$(ROLLBAR_TOKEN) -s -w"

default: build

all: clean windows linux macos

define build-os
	$(foreach GOARCH, $(ARCHITECTURES), \
		$(shell export GOOS=$(1); export GOARCH=$(GOARCH); go build ${LDFLAGS} -o output/$(BINARY)-$(1)-$(GOARCH)$(2); mkdir -p gh-pages/static/$(1)/$(GOARCH); cp output/$(BINARY)-$(1)-$(GOARCH)$(2) gh-pages/static/$(1)/$(GOARCH)/$(BINARY)$(2)))
endef

windows: 
	$(call build-os,windows,.exe)

linux: 
	$(call build-os,linux,)

macos: 
	$(call build-os,darwin,)

build:
	go build ${LDFLAGS}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete
