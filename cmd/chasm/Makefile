chasm: chasm.go main.go ../../pkg/vm/opcodes.go
	go build

chasm.go: chasm.peggo
	pigeon -o chasm.go chasm.peggo

clean:
	rm chasm chasm.go

test: chasm
	go test -v --race

.PHONY: examples
examples: chasm
	./chasm --output examples/quadratic.chbin --comment "Test of quadratic" examples/quadratic.chasm
	./chasm --output examples/majority.chbin --comment "Test of majority" examples/majority.chasm
