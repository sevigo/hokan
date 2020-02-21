BASEDIR=${CURDIR}
TMP=${BASEDIR}/tmp
VENDOR_TMP=${TMP}/vendor
LOCAL_BIN:=${TMP}/bin

run:
	cd cmd/hokan; go run main.go inject_server.go  inject_store.go  wire_gen.go

wire:
	go get github.com/google/wire/cmd/wire

generate: wire
	cd cmd/hokan && go generate	

install-mockgen:
	GOPATH=${TMP} go get github.com/golang/mock/gomock
	GOPATH=${TMP} go install github.com/golang/mock/mockgen

mockgen: install-mockgen
	${LOCAL_BIN}/mockgen -destination=mocks/mock_gen.go -package=mocks github.com/sevigo/hokan/pkg/core DirectoryStore,EventCreator

install-golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b ${LOCAL_BIN}

lint: install-golangci-lint
	${LOCAL_BIN}/golangci-lint run ./...