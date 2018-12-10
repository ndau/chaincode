package chain

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
)

func TestToValueScalar(t *testing.T) {
	tt, _ := time.Parse("2006-01-02T15:04:05Z", "2018-02-03T12:34:56Z")
	ts, _ := vm.NewTimestampFromTime(tt)
	tests := []struct {
		name    string
		args    interface{}
		want    vm.Value
		wantErr bool
	}{
		{"int", int(1), vm.NewNumber(1), false},
		{"int64", int64(1000), vm.NewNumber(1000), false},
		{"minus 1", int(-1), vm.NewNumber(-1), false},
		{"uint64", uint64(1), vm.NewNumber(1), false},
		{"illegal uint64", uint64(math.MaxUint64), nil, true},
		{"string", "hello", vm.NewBytes([]byte("hello")), false},
		{"time", tt, ts, false},
		{"true", true, vm.NewTrue(), false},
		{"false", false, vm.NewFalse(), false},
		{"[]int", []int{1, 23}, nil, true},
		{"map", map[int]int{1: 2}, nil, true},
		{"ptr to time", &tt, ts, false},
		{"timestamp", types.Timestamp(ts.T()), ts, false},
		{"undecorated struct", struct{ X int }{3}, vm.NewStruct(), false},
		{"unexpected type", complex64(0 - 1i), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToValueScalar(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToValueScalar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToValueScalar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToValue(t *testing.T) {
	// for time tests
	tt, _ := time.Parse("2006-01-02T15:04:05Z", "2018-02-03T12:34:56Z")
	ts, _ := vm.NewTimestampFromTime(tt)

	// for struct tests
	type st struct {
		H string `chain:"1"`
		J string `chain:"3"`
		N int    `chain:"2"`
		X string
		Y string `chain:"0"`
	}

	type foo int64
	type bar int64
	type custom struct {
		A int64 `chain:"0"`
		B foo   `chain:"1"`
		C bar   `chain:"2"`
	}

	type nest2 struct {
		P string `chain:"15"`
		Q string `chain:"16"`
	}

	type nest1 struct {
		H string `chain:"1"`
		J string `chain:"3"`
		M nest2  `chain:"."`
	}

	type nestp struct {
		H string `chain:"1"`
		J string `chain:"3"`
		M *nest2 `chain:"."`
	}

	type nonest1 struct {
		H string `chain:"1"`
		J string `chain:"3"`
		M nest2
	}

	type args struct {
		x interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    vm.Value
		wantErr bool
	}{
		{"simple", args{
			struct {
				X int `chain:"0"`
			}{3},
		}, vm.NewTestStruct(vm.NewNumber(3)), false},
		{"badtag", args{
			struct {
				X int `chain:"x"`
			}{3},
		}, nil, true},
		{"in order", args{
			struct {
				X int `chain:"0"`
				Y int `chain:"1"`
				Z int `chain:"2"`
			}{3, 4, 5},
		}, vm.NewTestStruct(vm.NewNumber(3), vm.NewNumber(4), vm.NewNumber(5)), false},
		{"out of order", args{
			struct {
				X int `chain:"2"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, vm.NewTestStruct(vm.NewNumber(4), vm.NewNumber(5), vm.NewNumber(3)), false},
		{"defined types", args{
			custom{3, 4, 5},
		}, vm.NewTestStruct(vm.NewNumber(3), vm.NewNumber(4), vm.NewNumber(5)), false},
		{"not continuous should not error", args{
			struct {
				X int `chain:"3"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, vm.NewStruct().Set(0, vm.NewNumber(4)).Set(1, vm.NewNumber(5)).Set(3, vm.NewNumber(3)), false},
		{"mixed types", args{
			struct {
				X string    `chain:"0"`
				Y int64     `chain:"1"`
				Z byte      `chain:"2"`
				T time.Time `chain:"3"`
			}{"hi", math.MaxInt64, 0x2A, tt},
		}, vm.NewTestStruct(vm.NewBytes([]byte("hi")), vm.NewNumber(math.MaxInt64), vm.NewNumber(42), ts), false},
		{"illegal field type", args{
			struct {
				X complex64 `chain:"0"`
			}{3},
		}, nil, true},
		{"simple array", args{
			[]int{1, 2, 3},
		}, vm.NewList().Append(vm.NewNumber(1)).Append(vm.NewNumber(2)).Append(vm.NewNumber(3)), false},
		{"struct array", args{
			[]st{
				st{"had", "job", 1, "nothing", "you"},
				st{"for", "phore", 4, "four", "fore"},
			},
		}, vm.NewList().Append(vm.NewTestStruct(
			vm.NewBytes([]byte("you")),
			vm.NewBytes([]byte("had")),
			vm.NewNumber(1),
			vm.NewBytes([]byte("job")),
		)).Append(vm.NewTestStruct(
			vm.NewBytes([]byte("fore")),
			vm.NewBytes([]byte("for")),
			vm.NewNumber(4),
			vm.NewBytes([]byte("phore")),
		)), false},
		{"nested struct", args{
			nest1{"a", "b", nest2{"c", "d"}}},
			vm.NewStruct().
				Set(1, vm.NewBytes([]byte("a"))).
				Set(3, vm.NewBytes([]byte("b"))).
				Set(15, vm.NewBytes([]byte("c"))).
				Set(16, vm.NewBytes([]byte("d"))),
			false},
		{"nested struct ptr", args{
			nestp{"a", "b", &nest2{"c", "d"}}},
			vm.NewStruct().
				Set(1, vm.NewBytes([]byte("a"))).
				Set(3, vm.NewBytes([]byte("b"))).
				Set(15, vm.NewBytes([]byte("c"))).
				Set(16, vm.NewBytes([]byte("d"))),
			false},
		{"nested struct with no . tag", args{
			nonest1{"a", "b", nest2{"c", "d"}}},
			vm.NewStruct().
				Set(1, vm.NewBytes([]byte("a"))).
				Set(3, vm.NewBytes([]byte("b"))),
			false},
		{"int", args{int(1)}, vm.NewNumber(1), false},
		{"int64", args{int64(1000)}, vm.NewNumber(1000), false},
		{"minus 1", args{int(-1)}, vm.NewNumber(-1), false},
		{"uint64", args{uint64(1)}, vm.NewNumber(1), false},
		{"illegal uint64", args{uint64(math.MaxUint64)}, nil, true},
		{"string", args{"hello"}, vm.NewBytes([]byte("hello")), false},
		{"[]byte", args{[]byte("hello")}, vm.NewBytes([]byte("hello")), false},
		{"time", args{tt}, ts, false},
		{"true", args{true}, vm.NewTrue(), false},
		{"false", args{false}, vm.NewFalse(), false},
		{"[]int", args{[]int{1, 23}}, vm.NewList().Append(vm.NewNumber(1)).Append(vm.NewNumber(23)), false},
		{"[] illegal values", args{[]complex64{1, 23i}}, nil, true},
		{"map", args{map[int]int{1: 2}}, nil, true},
		{"ptr to time", args{&tt}, ts, false},
		{"[][]int", args{[][]int{[]int{1}, []int{2, 3}}},
			vm.NewList().Append(vm.NewList().Append(vm.NewNumber(1))).Append(vm.NewList().Append(vm.NewNumber(2)).Append(vm.NewNumber(3))), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToValue(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractConstants(t *testing.T) {
	type nest2 struct {
		P string `chain:"15"`
		Y int64
		Q string `chain:"16"`
	}

	type nest1 struct {
		H string `chain:"1"`
		X int64
		J string `chain:"3"`
		M nest2  `chain:"."`
	}

	type nestp struct {
		H string `chain:"1"`
		J string `chain:"3"`
		M *nest2 `chain:"."`
	}

	type nonest1 struct {
		H string `chain:"1"`
		J string `chain:"3"`
		M nest2
	}

	type args struct {
		x interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]byte
		wantErr bool
	}{
		{"simple", args{
			struct {
				X int `chain:"0"`
			}{3},
		}, map[string]byte{"X": 0}, false},
		{"in order", args{
			struct {
				X int `chain:"0"`
				Y int `chain:"1"`
				Z int `chain:"2"`
			}{3, 4, 5},
		}, map[string]byte{"X": 0, "Y": 1, "Z": 2}, false},
		{"out of order", args{
			struct {
				X int `chain:"2"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, map[string]byte{"X": 2, "Y": 0, "Z": 1}, false},
		{"not continuous should not error", args{
			struct {
				X int `chain:"3"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, map[string]byte{"X": 3, "Y": 0, "Z": 1}, false},
		{"mixed types", args{
			struct {
				X string `chain:"0"`
				Y int64  `chain:"1"`
				Z byte   `chain:"2"`
			}{"hi", math.MaxInt64, 0x2A},
		}, map[string]byte{"X": 0, "Y": 1, "Z": 2}, false},
		{"rename", args{
			struct {
				X int `chain:"0,foo"`
				Y int `chain:"1,Bar"`
				Z int `chain:"2"`
			}{3, 4, 5},
		}, map[string]byte{"FOO": 0, "BAR": 1, "Z": 2}, false},
		{"bad number", args{
			struct {
				X int `chain:"x"`
			}{3},
		}, nil, true},
		{"illegal name", args{
			struct {
				X int `chain:"0,a+b"`
			}{3},
		}, nil, true},
		{"empty name", args{
			struct {
				X int `chain:"0,"`
			}{3},
		}, nil, true},
		{"only some fields", args{
			struct {
				X int `chain:"0"`
				Y int
				Z int `chain:"1"`
			}{3, 4, 5},
		}, map[string]byte{"X": 0, "Z": 1}, false},
		{"no chain tags", args{
			struct {
				X int
				Y int
				Z int
			}{3, 4, 5},
		}, nil, true},
		{"not a struct", args{[]int{3}}, nil, true},
		{"nested struct", args{
			nest1{"a", 0, "b", nest2{"c", 0, "d"}}},
			map[string]byte{"H": 1, "J": 3, "P": 15, "Q": 16}, false},
		{"nested struct with no . tag", args{
			nonest1{"a", "b", nest2{"c", 0, "d"}}},
			map[string]byte{"H": 1, "J": 3}, false},
		{"nested struct with ptr", args{
			nestp{"a", "b", &nest2{"c", 0, "d"}}},
			map[string]byte{"H": 1, "J": 3, "P": 15, "Q": 16}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractConstants(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractConstants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractConstants() = %v, want %v", got, tt.want)
			}
		})
	}
}
