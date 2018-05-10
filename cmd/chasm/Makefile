chasm: chasm.go main.go
	go build

chasm.go: chasm.peggo
	pigeon -o chasm.go chasm.peggo

clean:
	rm chasm chasm.go

test: chasm
	go test -v --race
