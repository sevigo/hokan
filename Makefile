
run:
	cd cmd/hokan; go run main.go inject_server.go  inject_store.go  wire_gen.go

wire:
	go get github.com/google/wire/cmd/wire

generate: wire
	cd cmd/hokan && go generate	