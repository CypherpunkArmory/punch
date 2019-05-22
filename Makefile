# ########################################################## #
# Includes cross-compiling, installation, cleanup
# ########################################################## #

SHELL=/bin/bash

# Check for required command tools to build or stop immediately
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=punch
VERSION=
ROLLBAR_TOKEN=
ARCHITECTURES=386 amd64
LINUX_ARCHITECTURES=386 amd64 arm arm64
ARM_SUB_ARCHITECTURES=5 6 7
# If more than one release is made in a day this should be incremented
# I did not find a good way to do this in CircleCI
CALVER=`date +%Y.%m.%d`.1
# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X github.com/cypherpunkarmory/punch/restapi.apiVersion=$(CALVER) -X github.com/cypherpunkarmory/punch/cmd.version=$(VERSION) -X github.com/cypherpunkarmory/punch/cmd.rollbarToken=$(ROLLBAR_TOKEN) -s -w"

default: build

release: all zip

all: clean windows linux macos

define build-os
	if [ "$(1)" = "linux" ]; then \
		for GOARCH in $(LINUX_ARCHITECTURES); do \
			if [ "$$GOARCH" = "arm" ]; then \
				for GOARM in $(ARM_SUB_ARCHITECTURES); do \
					echo "building $(1) $${GOARCH}v$${GOARM}"; \
					export GOOS=$(1); export GOARCH=$$GOARCH; export GOARM=$$GOARM; go build ${LDFLAGS} -o output/$(BINARY)-$(1)-$${GOARCH}v$${GOARM}$(2); \
				done; \
			else \
				echo "building $(1) $$GOARCH"; \
				export GOOS=$(1); export GOARCH=$$GOARCH; go build ${LDFLAGS} -o output/$(BINARY)-$(1)-$${GOARCH}$(2); \
			fi; \
		done; \
	else \
		for GOARCH in $(ARCHITECTURES); do \
			echo "building $(1) $$GOARCH"; \
			export GOOS=$(1); export GOARCH=$$GOARCH; go build ${LDFLAGS} -o output/$(BINARY)-$(1)-$${GOARCH}$(2); \
		done; \
	fi;
endef

windows: 
	$(call build-os,windows,.exe)

linux: 
	$(call build-os,linux,)

macos: 
	$(call build-os,darwin,)

build:
	go build ${LDFLAGS}

zip:
	mkdir -p output/release; \
	cd output; \
	echo "export const PunchVersion = '$(VERSION)';" > release/version.js; \
	for f in punch*;  do \
		case "$$f" in \
			*.exe) extension=".exe" ;; \
			*) extension="" ;; \
		esac; \
		filename="$${f%.*}"; \
		cp $${f} punch$${extension}; \
		zip release/$${filename}.zip punch$${extension};  \
		rm punch$${extension}; \
	done;

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete
