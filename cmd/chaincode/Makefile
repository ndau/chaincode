chaincode: chaincode.go main.go ../../pkg/vm/opcodes.go
	go build

chaincode.go: chaincode.peggo
	pigeon -o chaincode.go chaincode.peggo

clean:
	rm chaincode chaincode.go

test: chaincode
	go test -v --race

.PHONY: examples
examples:
	./chaincode --output examples/quadratic.chbin --comment "Test of quadratic" examples/quadratic.chaincode
