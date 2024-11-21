GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
ifeq (${CODE},)
	CODE := Stew
else
	CODE := ${CODE}
endif
PROGRAM := $(CODE)
ifeq ($(GOOS),windows)
	PROGRAM := $(PROGRAM).exe
	GOFLAGS := -v -buildmode=exe
else
	GOFLAGS := -v
endif

CGO_ENABLED := 0
COMMIT_HASH := $(shell git diff --quiet || echo "local")
ifeq ($(COMMIT_HASH),)
	COMMIT_HASH := $(shell git rev-parse --short=10 HEAD)
endif

LDFLAGS := -X stew/embeds.ExecutableVersion=$(COMMIT_HASH)

LDFLAGS := $(LDFLAGS) -X stew/embeds.Code=$(CODE)

RELEASE_MODE := ${RELEASE_MODE}
ifeq ($(RELEASE_MODE),)
	RELEASE_MODE := release
endif

DOCKER_TAG_CODE := ${DOCKER_TAG_CODE}
ifeq ($(DOCKER_TAG_CODE),)
	DOCKER_TAG_CODE := "$(CODE):local"
endif
DOCKER_TAG_CODE := $(shell echo "$(DOCKER_TAG_CODE)" | tr '[:upper:]' '[:lower:]')

all: test build

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags "$(LDFLAGS) -X stew/router.ginMode=$(RELEASE_MODE)" -o $(PROGRAM)

.PHONY: test
test:
	go test -ldflags "$(LDFLAGS) -X stew/router.ginMode=debug" -count=1 -v ./test/...

.PHONY: clean
clean:
	rm -f $(PROGRAM)