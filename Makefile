.PHONY: generate clean fuzz fuzzmillion benchmarks test examples all

all: clean generate test fuzz benchmarks examples coverage

fuzz: test
	FUZZ_RUNS=50000 go test --race -v -timeout 1m ./pkg/vm -run "TestFuzz*" -coverprofile=/tmp/coverfuzz

fuzzmillion: test
	FUZZ_RUNS=1000000 go test --race -v -timeout 2h ./pkg/vm -run "TestFuzz*" -coverprofile=/tmp/coverfuzz

benchmarks:
	go test -bench ./pkg/vm -benchmem

generate: opcodes.md pkg/vm/opcodes.go pkg/vm/miniasmOpcodes.go pkg/vm/opcode_string.go \
		pkg/vm/extrabytes.go cmd/chasm/chasm.peggo pkg/vm/enabledopcodes.go

clean:
	rm -f cmd/opcodes/opcodes
	rm -f cmd/chasm/chasm cmd/chasm/chasm.go

opcodes.md: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --opcodes opcodes.md

pkg/vm/opcodes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --defs pkg/vm/opcodes.go

pkg/vm/miniasmOpcodes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --miniasm pkg/vm/miniasmOpcodes.go

pkg/vm/extrabytes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --extra pkg/vm/extrabytes.go

pkg/vm/enabledopcodes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --enabled pkg/vm/enabledopcodes.go

cmd/chasm/chasm.peggo: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --pigeon cmd/chasm/chasm.peggo

cmd/opcodes/opcodes: cmd/opcodes/*.go
	cd cmd/opcodes && go build

pkg/vm/opcode_string.go: pkg/vm/opcodes.go
	go generate ./pkg/vm

chasm: cmd/chasm/chasm.go pkg/vm/opcodes.go cmd/chasm/*.go
	go build -o ./cmd/chasm/chasm ./cmd/chasm

cmd/chasm/chasm.go: cmd/chasm/chasm.peggo
	pigeon -o ./cmd/chasm/chasm.go ./cmd/chasm/chasm.peggo

test: cmd/chasm/chasm.go pkg/vm/*.go pkg/chain/*.go chasm
	rm -f /tmp/cover*
	go test ./pkg/chain -v --race -timeout 10s -coverprofile=/tmp/coverchain
	go test ./cmd/chasm -v --race -timeout 10s -coverprofile=/tmp/coverchasm
	go test ./pkg/vm -v --race -timeout 10s -coverprofile=/tmp/covervm

examples: chasm
	./cmd/chasm/chasm --output cmd/chasm/examples/quadratic.chbin --comment "Test of quadratic" cmd/chasm/examples/quadratic.chasm
	./cmd/chasm/chasm --output cmd/chasm/examples/majority.chbin --comment "Test of majority" cmd/chasm/examples/majority.chasm
	./cmd/chasm/chasm --output cmd/chasm/examples/onePlus1of3.chbin --comment "1+1of3" cmd/chasm/examples/onePlus1of3.chasm

