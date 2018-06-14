package vm

import (
	"math"
	"reflect"
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

func TestToBytes(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want []byte
	}{
		{"a", 1, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"b", 0x010203, []byte{0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"c", -1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{"d", -256, []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{"e", int64(math.MaxInt64), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}},
		{"f", int64(math.MinInt64), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
		{"g", 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBytes(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mulDiv(t *testing.T) {
	type args struct {
		v int64
		n int64
		d int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"a", args{100, 5, 8}, 62, false},
		{"b", args{123456, 567890, 7890987}, 8884, false},
		{"c", args{math.MaxInt64, 17, 23}, 6817274983762225596, false},
		{"d", args{-123456, 567890, 7890987}, -8884, false},
		{"e", args{100, 0, 8}, 0, false},
		{"f", args{100, 20, 8}, 250, false},
		{"g", args{100, 5, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mulDiv(tt.args.v, tt.args.n, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("mulDiv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mulDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_divmod(t *testing.T) {
	type args struct {
		n int64
		d int64
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{"a", args{17, 5}, 3, 2},
		{"b", args{15, 5}, 3, 0},
		{"c", args{15, 25}, 0, 15},
		{"d", args{2, 3}, 0, 2},
		{"e", args{-17, 5}, -3, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := divmod(tt.args.n, tt.args.d)
			if got != tt.want {
				t.Errorf("divmod() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("divmod() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
