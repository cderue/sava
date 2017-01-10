# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR:
.SUFFIXES:

# Constants, these can be overwritten in your Makefile.local
CONTAINER := magneticio/buildserver:0.4

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

# Targets
.PHONY: all
all: default

# Using our buildserver which contains all the necessary dependencies
.PHONY: default
default:
	docker run \
		--interactive \
		--rm \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume $(shell command -v docker):/usr/bin/docker \
		--volume $(CURDIR):/srv/src/go/src/github.com/magneticio/sava \
		--workdir=/srv/src/go/src/github.com/magneticio/sava \
		$(CONTAINER) \
			make build


.PHONY: build
build:
	$(CURDIR)/docker/build.sh 0


.PHONY: push
push:
	$(CURDIR)/docker/push.sh 0

.PHONY: clean
clean:
	rm -rf $(CURDIR)/docker/target/

