OUT_DIR := ./build
OUT := $(OUT_DIR)/gstore
MAIN_PKG := ./cmd/gstore

GOBIN := $(shell go env GOPATH)/bin

MOCKGEN_NAME := mockgen
MOCKGEN_VERSION := v1.4.4

default: build

.PHONY: build
build:
	CGO_ENABLED=0 go build -mod=vendor -o $(OUT) $(MAIN_PKG)

.PHONY: linux
linux: export GOOS := linux
linux: export GOARCH := amd64
linux: LINUX_OUT := $(OUT)-$(GOOS)-$(GOARCH)
linux:
	@echo BUILDING $(LINUX_OUT)
	CGO_ENABLED=0 go build -mod=vendor -o $(LINUX_OUT) $(MAIN_PKG)
	@echo DONE

.PHONY: image
image:
	docker build -t gstore-local -f scripts/Dockerfile .

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

.PHONY: test
test: GO_TEST_FLAGS := -race
test:
	go test -mod=vendor -v $(GO_TEST_FLAGS) $(GO_TEST_TAGS) ./...

.PHONY: fulltest
fulltest: GO_TEST_TAGS := --tags=integration
fulltest: test

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: generate
generate: check-all
	go generate -mod=vendor -x ./...

.PHONY: install-mockgen
install-mockgen: MOCKGEN_INSTALL_PATH := $(GOBIN)/$(MOCKGEN_NAME)
install-mockgen:
	@echo INSTALLING $(MOCKGEN_INSTALL_PATH) $(MOCKGEN_NAME)
	# we need to change dir to use go modules without updating repo deps
	cd $(TMPDIR) && GO111MODULE=on go get github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
	@echo DONE

.PHONY: check-mockgen-version
check-mockgen-version: ACTUAL_MOCKGEN_VERSION := $(shell $(MOCKGEN_NAME) --version 2>/dev/null)
check-mockgen-version:
	[ -z $(ACTUAL_MOCKGEN_VERSION) ] && \
		echo 'Mockgen is not installed, run `make mockgen-install`' && \
		exit 1 || true

	if [ $(ACTUAL_MOCKGEN_VERSION) != $(MOCKGEN_VERSION) ] ; then \
		echo $(MOCKGEN_NAME) is version $(ACTUAL_MOCKGEN_VERSION), want $(MOCKGEN_VERSION) ; \
		echo 'Make sure $$GOBIN has precedence in $$PATH and' \
		'run `make mockgen-install` to install the correct version' ; \
    exit 1 ; \
	fi

.PHONY: check-all
check-all: check-mockgen-version

.PHONY: install-all
install-all: install-mockgen

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir scripts/migrations/postgres -seq $(NAME)