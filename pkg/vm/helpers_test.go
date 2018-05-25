package vm

import (
	"math"
	"testing"
)

func TestFractionLess(t *testing.T) {
	type args struct {
		n1 int64
		d1 int64
		n2 int64
		d2 int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{3, 10, 2, 10}, false},
		{"test2", args{2, 10, 3, 10}, true},
		{"test3", args{2, 10, 2, 10}, false},
		{"test4", args{2, 10, 21, 100}, true},
		{"test5", args{2, 10, 19, 100}, false},
		{"test6", args{3, 10, 29, 100}, false},
		{"test7", args{19, 100, 195, 1000}, true},
		{"testpi", args{355, 113, 22, 7}, true},
		// .3141593 vs .314158655
		{"test8", args{3141593, 10000000, 1349302, 4294970}, false},
		{"testBig1", args{math.MaxInt64 / 3, math.MaxInt64, 1, 4}, false},
		{"testBig2", args{math.MaxInt64 / 3, math.MaxInt64, 1, 2}, true},
		{"testBig3", args{math.MaxInt64 / 2, math.MaxInt64, 0, 15}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FractionLess(tt.args.n1, tt.args.d1, tt.args.n2, tt.args.d2); got != tt.want {
				t.Errorf("FractionLess() = %v, want %v", got, tt.want)
			}
		})
	}
}
