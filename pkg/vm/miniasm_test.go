package vm

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

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
