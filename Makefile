USERSPACE=Confialink
NAMESPACE=wallet
SERVICE=users

APP=./build/service_${SERVICE}
PROJECT?=github.com/${USERSPACE}/${NAMESPACE}-${SERVICE}
DATE := $(shell date +'%Y.%m.%d %H:%M:%S')

ndef = $(if $(value $(1)),,$(error required environment variable $(1) is not set))

ifndef COMMIT
	COMMIT := $(shell git rev-parse HEAD)
endif

ifndef TAG
	TAG = $(shell git describe --exact-match --tags $(git log -n1 --pretty='%h') 2>/dev/null)
endif

GOOS?=linux
DOCKER_TAG?=wallet-users
GO111MODULE?=on
GOPRIVATE?=github.com/Confialink

show:
	@echo ${PROJECT}
	@echo ${DATE}

fast_build:
	CGO_ENABLED=0 GOOS=${GOOS} go build \
		-gcflags "all=-N -l" \
		-ldflags '-X "${PROJECT}/internal/version.DATE=${DATE}" -X ${PROJECT}/internal/version.COMMIT=${COMMIT} -X ${PROJECT}/internal/version.TAG=${TAG}' \
		-o ${APP} ./cmd/.

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo \
		-ldflags '-s -w -X "${PROJECT}/internal/version.DATE=${DATE}" -X ${PROJECT}/internal/version.COMMIT=${COMMIT} -X ${PROJECT}/internal/version.TAG=${TAG}' \
		-o ${APP} ./cmd/.

docker-build:
	$(call ndef,REPOSITORY_PRIVATE_KEY)
	docker build . --build-arg REPOSITORY_PRIVATE_KEY --build-arg TAG=${TAG} -t ${DOCKER_TAG}

gen-protobuf:
	protoc --proto_path=. --go_out=. --twirp_out=. rpc/proto/users/users.proto

clean:
	@[ -f ${APP} ] && rm -f ${APP} || true
