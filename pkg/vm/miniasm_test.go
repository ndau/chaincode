package vm

import (
	"testing"
)

// just ensure that MiniAsmSafe does not in fact panic
func TestMiniAsmSafe(t *testing.T) {
	tests := []struct {
		arg     string
		wantErr bool
	}{
		{"This is definitely not valid chaincode.", true},
		{"ff", false},
		{"handler 0 zero enddef", false},
		{"ffoo", true},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			_, err := MiniAsmSafe(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("MiniAsmSafe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
