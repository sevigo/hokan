
run:
	cd cmd/hokan; go run .\main.go .\wire_gen.go .\inject_server.go 

wire:
	go get github.com/google/wire/cmd/wire

generate: wire
	cd cmd/hokan && go generate	