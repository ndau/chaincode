CHASM = cmd/chasm/chasm
CHAIN = cmd/chain/chain
CRANK = cmd/crank/crank
CHFMT = cmd/chfmt/chfmt
EXAMPLES = cmd/chasm/examples
OPCODES = cmd/opcodes/opcodes

.PHONY: generate clean fuzz fuzzmillion benchmarks \
	test examples all build chasm crank chfmt opcodes format

all: clean generate build test fuzz benchmarks format examples coverage

build: generate opcodes chasm crank chfmt

opcodes: $(OPCODES)

crank: $(CRANK)

chasm: $(CHASM)

chfmt: $(CHFMT)

fuzz: test
	FUZZ_RUNS=10000 go test --race -v -timeout 1m ./pkg/vm -run "TestFuzz*" -coverprofile=/tmp/coverfuzz

fuzzmillion: test
	FUZZ_RUNS=1000000 go test --race -v -timeout 2h ./pkg/vm -run "TestFuzz*" -coverprofile=/tmp/coverfuzz

benchmarks:
	cd pkg/vm && go test -bench=. -benchmem

generate: opcodes.md pkg/vm/opcodes.go pkg/vm/miniasmOpcodes.go pkg/vm/opcode_string.go \
		pkg/vm/extrabytes.go pkg/vm/enabledopcodes.go \
		cmd/chasm/chasm.peggo cmd/chasm/predefined.go

clean:
	rm -f $(OPCODES)
	rm -f $(CHASM) cmd/chasm/chasm.go
	rm -f $(CRANK)
	rm -f $(CHFMT) cmd/chfmt/chfmt.go

opcodes.md: opcodes
	$(OPCODES) --opcodes opcodes.md

pkg/vm/opcodes.go: opcodes
	$(OPCODES) --defs pkg/vm/opcodes.go

pkg/vm/miniasmOpcodes.go: opcodes
	$(OPCODES) --miniasm pkg/vm/miniasmOpcodes.go

pkg/vm/extrabytes.go: opcodes
	$(OPCODES) --extra pkg/vm/extrabytes.go

pkg/vm/enabledopcodes.go: opcodes
	$(OPCODES) --enabled pkg/vm/enabledopcodes.go

cmd/chasm/chasm.peggo: opcodes
	$(OPCODES) --pigeon cmd/chasm/chasm.peggo

cmd/chasm/predefined.go: opcodes
	$(OPCODES) --consts cmd/chasm/predefined.go

$(OPCODES): cmd/opcodes/*.go
	cd cmd/opcodes && go build

pkg/vm/opcode_string.go: pkg/vm/opcodes.go
	go generate ./pkg/vm

$(CHASM): cmd/chasm/chasm.go pkg/vm/opcodes.go cmd/chasm/*.go
	go build -o $(CHASM) ./cmd/chasm

cmd/chasm/chasm.go: cmd/chasm/chasm.peggo
	pigeon -o ./cmd/chasm/chasm.go ./cmd/chasm/chasm.peggo

test: cmd/chasm/chasm.go pkg/vm/*.go pkg/chain/*.go chasm
	rm -f /tmp/cover*
	go test ./pkg/chain -v --race -timeout 10s -coverprofile=/tmp/coverchain
	go test ./cmd/chasm -v --race -timeout 10s -coverprofile=/tmp/coverchasm
	go test ./pkg/vm -v --race -timeout 10s -coverprofile=/tmp/covervm

examples: chasm
	$(CHASM) --output $(EXAMPLES)/quadratic.chbin --comment "Test of quadratic" $(EXAMPLES)/quadratic.chasm
	$(CHASM) --output $(EXAMPLES)/majority.chbin --comment "Test of majority" $(EXAMPLES)/majority.chasm
	$(CHASM) --output $(EXAMPLES)/onePlus1of3.chbin --comment "1+1of3" $(EXAMPLES)/onePlus1of3.chasm
	$(CHASM) --output $(EXAMPLES)/first.chbin --comment "the first key must be set" $(EXAMPLES)/first.chasm
	$(CHASM) --output $(EXAMPLES)/one.chbin --comment "unconditionally return numeric 1" $(EXAMPLES)/one.chasm
	$(CHASM) --output $(EXAMPLES)/zero.chbin --comment "returns numeric 0 in all cases" $(EXAMPLES)/zero.chasm
	$(CHASM) --output $(EXAMPLES)/rfe.chbin --comment "standard RFE rules" $(EXAMPLES)/rfe.chasm

format: chfmt
	$(CHFMT) -O $(EXAMPLES)/quadratic.chasm
	$(CHFMT) -O $(EXAMPLES)/majority.chasm
	$(CHFMT) -O $(EXAMPLES)/onePlus1of3.chasm
	$(CHFMT) -O $(EXAMPLES)/first.chasm
	$(CHFMT) -O $(EXAMPLES)/one.chasm
	$(CHFMT) -O $(EXAMPLES)/zero.chasm
	$(CHFMT) -O $(EXAMPLES)/rfe.chasm

$(CRANK): cmd/crank/*.go cmd/crank/glide.* generate
	go build -o $(CRANK) ./cmd/crank

cmd/chfmt/chfmt.go: cmd/chfmt/chfmt.peggo
	pigeon -o ./cmd/chfmt/chfmt.go ./cmd/chfmt/chfmt.peggo

$(CHFMT): cmd/chfmt/*.go cmd/chfmt/chfmt.go
	go build -o $(CHFMT) ./cmd/chfmt

