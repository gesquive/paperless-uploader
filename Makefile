#
#  Makefile
#
#  A kickass golang v1.20.x makefile
#  v1.20.0

export SHELL ?= /bin/bash
include make.cfg

GOCC := go

# Program version
MK_VERSION := $(shell git describe --always --tags --dirty)
MK_HASH := $(shell git rev-parse --short HEAD)
MK_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

PKG_NAME := ${REPO_HOST_URL}/${OWNER}/${PROJECT_NAME}
PKG_PATH := ${GOPATH}/src/${PKG_NAME}
PKG_LIST := ./...

COVER_PATH := coverage
DIST_PATH ?= dist
INSTALL_PATH ?= /usr/local/bin/

CF_NAME := ${REGISTRY_URL}/${OWNER}/${PROJECT_NAME}
CF_VERSION = $(shell git describe --always --tags | sed 's/^v//' | sed 's/-g/-/')
CF_PLATFORMS ?= linux/amd64,linux/arm/v7,linux/arm64
CF_PATH ?= Containerfile

BIN ?= ${GOPATH}/bin
CHECK ?= ${BIN}/staticcheck
GORELEASER ?= ${BIN}/goreleaser
DOCKER ?= docker

export CGO_ENABLED = 0
export DOCKER_CLI_EXPERIMENTAL = enabled

default: test build

.PHONY: help
help:
	@echo 'Management commands for $(PROJECT_NAME):'
	@grep -Eh '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	 awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compile the project
	@echo "building ${OWNER} ${BIN_NAME} ${MK_VERSION}"
	@echo "GOPATH=${GOPATH}"
	${GOCC} build -a -ldflags "-X main.buildVersion=${MK_VERSION} -X main.buildDate=${MK_DATE} -X main.buildCommit=${MK_HASH}" -o ${BIN_NAME}

.PHONY: install
install: build ## Install the binary
	install -d ${INSTALL_PATH}
	install -m 755 ./${BIN_NAME} ${INSTALL_PATH}/${BIN_NAME}

.PHONY: link
link: $(PKG_PATH) ## Symlink this project into the GOPATH
$(PKG_PATH):
	@mkdir -p `dirname $(PKG_PATH)`
	@ln -s $(PWD) $(PKG_PATH) >/dev/null 2>&1

.PHONY: path # Returns the project path
path:
	@echo $(PKG_PATH)

.PHONY: deps
deps: ## Download project dependencies
	${GOCC} mod download
	${GOCC} mod verify

.PHONY: check
check: ${CHECK} ## Lint/Check the source code
	${CHECK} ${PKG_LIST}

.PHONY: test
test: ## Run golang tests
	${GOCC} test ${PKG_LIST}

.PHONY: bench
bench: ## Run golang benchmarks
	${GOCC} test -benchmem -bench=. ${PKG_LIST}

.PHONY: coverage
coverage: ## Run coverage report
	${GOCC} test -v -cover ${PKG_LIST}

.PHONY: coverage-report
coverage-report: ## Generate global code coverage report
	mkdir -p "${COVER_PATH}"
	${GOCC} test -v -coverprofile "${COVER_PATH}/coverage.dat" ${PKG_LIST}
	${GOCC} tool cover -html="${COVER_PATH}/coverage.dat" -o "${COVER_PATH}/coverage.html"

.PHONY: race
race: ## Run data race detector
	CGO_ENABLED=1 ${GOCC} test -race ${PKG_LIST}

.PHONY: clean
clean: ## Clean the directory tree
	${GOCC} clean
	rm -f ./${BIN_NAME}.test
	rm -f ./${BIN_NAME}
	rm -rf "${DIST_PATH}"
	rm -f "${COVER_PATH}"

.PHONY: release-snapshot
release-snapshot: ${GORELEASER} ## Cross compile and package to local disk
	${GORELEASER} release --skip=publish --clean --snapshot

.PHONY: release
release: ${GORELEASER} ## Cross compile and package the full distribution
	${GORELEASER} release --clean

.PHONY: fmt
fmt: ## Reformat the source tree with gofmt
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

# Install golang dependencies here
${BIN}/%: 
	@echo "Installing ${PACKAGE} to ${BIN}"
	@mkdir -p ${BIN}
	@tmp=$$(mktemp -d); \
       env GOPATH=$$tmp GOBIN=${BIN} ${GOCC} install ${PACKAGE} \
        || ret=$$?; \
       rm -rf $$tmp ; exit $$ret

${BIN}/golint:     PACKAGE=honnef.co/go/tools/cmd/staticcheck@2023.1.6
${BIN}/goreleaser: PACKAGE=github.com/goreleaser/goreleaser@v1.23.0

# Image related targets
.PHONY: build-docker
build-docker: ## Build the image
	@echo "building ${MK_VERSION}"
	${DOCKER} info
	${DOCKER} build -f ${CF_PATH} --build-arg TARGETARCH=amd64 --build-arg TARGETOS=linux --pull -t ${CF_NAME}:${MK_VERSION} .

# build manifest for git describe
# manifest version is "1.2.3-g23ab3df"
# image version is "1.2.3-g23ab3df-amd64"

.PHONY: init-image-build
init-image-build:
	${DOCKER} context create build
	${DOCKER} buildx create --driver docker-container --name gobuild --use build
	${DOCKER} buildx inspect --bootstrap
	${DOCKER} buildx ls

.PHONY: release-image-snapshot
release-image-snapshot: init-image-build
	@echo "building multi-arch image ${CF_VERSION}"
	${DOCKER} buildx build -f ${CF_PATH} --platform ${CF_PLATFORMS} --pull -t ${CF_NAME}:${CF_VERSION} --push .

.PHONY: release-docker
release-image: init-image-build ## Build a multi-arch manifest and images
	@echo "building multi-arch image ${CF_VERSION}"
	${DOCKER} buildx build -f ${CF_PATH} --platform ${CF_PLATFORMS} --pull -t ${CF_NAME}:${CF_VERSION} -t ${CF_NAME}:latest --push .
