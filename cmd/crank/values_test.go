package main

import (
	"reflect"
	"testing"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
)

func Test_parseValues(t *testing.T) {
	loadConstants()

	tests := []struct {
		name    string
		s       string
		want    []vm.Value
		wantErr bool
	}{
		{"number", "123", []vm.Value{vm.NewNumber(123)}, false},
		{"struct", "{ ACCT_BALANCE: 12345 }",
			[]vm.Value{vm.NewStruct().Set(61, vm.NewNumber(12345))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValues(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
