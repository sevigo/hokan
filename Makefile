BASEDIR=${CURDIR}
TMP=${BASEDIR}/tmp
VENDOR_TMP=${TMP}/vendor
LOCAL_BIN:=${TMP}/bin

run:
	cd cmd/hokan; go run main.go application.go inject_server.go inject_store.go inject_watcher.go inject_backup.go wire_gen.go

build:
	cd cmd/hokan; go build -o ../../bin/hokan; cd ../..; ./bin/hokan

wire:
	go get github.com/google/wire/cmd/wire

# run this command after changing something in cmd/hokan/inject_*
generate: wire
	go generate	./...

install-mockgen:
	go get github.com/golang/mock/mockgen@v1.5.0

mockgen: install-mockgen
	# grep "interface {" pkg/core/* | awk '{print $2}' | paste -sd "," - 
	mockgen -destination=mocks/mock_gen.go -package=mocks github.com/sevigo/hokan/pkg/core Backup,DB,DirectoryStore,EventCreator,FileStore,MinioWrapper,Notifier,ServerSideEventCreator,UserStore,Watcher

install-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${LOCAL_BIN} v1.24.0

lint: install-golangci-lint
	${LOCAL_BIN}/golangci-lint run

test:
	go test -timeout 10s -cover ./...
