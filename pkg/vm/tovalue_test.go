package vm

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestToValueScalar(t *testing.T) {
	tt, _ := time.Parse("2006-01-02T15:04:05Z", "2018-02-03T12:34:56Z")
	ts, _ := NewTimestampFromTime(tt)
	tests := []struct {
		name    string
		args    interface{}
		want    Value
		wantErr bool
	}{
		{"int", int(1), NewNumber(1), false},
		{"int64", int64(1000), NewNumber(1000), false},
		{"minus 1", int(-1), NewNumber(-1), false},
		{"uint64", uint64(1), NewNumber(1), false},
		{"illegal uint64", uint64(math.MaxUint64), nil, true},
		{"string", "hello", NewBytes([]byte("hello")), false},
		{"time", tt, ts, false},
		{"true", true, NewNumber(1), false},
		{"false", false, NewNumber(0), false},
		{"[]int", []int{1, 23}, nil, true},
		{"map", map[int]int{1: 2}, nil, true},
		{"ptr to time", &tt, ts, false},
		{"illegal struct", struct{ X int }{3}, nil, true},
		{"unexpected type", int32(17), nil, true},
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
	ts, _ := NewTimestampFromTime(tt)

	// for struct tests
	type st struct {
		H string `chain:"1"`
		J string `chain:"3"`
		N int    `chain:"2"`
		X string
		Y string `chain:"0"`
	}

	type args struct {
		x interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Value
		wantErr bool
	}{
		{"simple", args{
			struct {
				X int `chain:"0"`
			}{3},
		}, NewTestStruct(NewNumber(3)), false},
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
		}, NewTestStruct(NewNumber(3), NewNumber(4), NewNumber(5)), false},
		{"out of order", args{
			struct {
				X int `chain:"2"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, NewTestStruct(NewNumber(4), NewNumber(5), NewNumber(3)), false},
		{"not continuous should not error", args{
			struct {
				X int `chain:"3"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, NewStruct().Set(0, NewNumber(4)).Set(1, NewNumber(5)).Set(3, NewNumber(3)), false},
		{"mixed types", args{
			struct {
				X string    `chain:"0"`
				Y int64     `chain:"1"`
				Z byte      `chain:"2"`
				T time.Time `chain:"3"`
			}{"hi", math.MaxInt64, 0x2A, tt},
		}, NewTestStruct(NewBytes([]byte("hi")), NewNumber(math.MaxInt64), NewNumber(42), ts), false},
		{"illegal field type", args{
			struct {
				X int32 `chain:"0"`
			}{3},
		}, nil, true},
		{"simple array", args{
			[]int{1, 2, 3},
		}, NewList().Append(NewNumber(1)).Append(NewNumber(2)).Append(NewNumber(3)), false},
		{"struct array", args{
			[]st{
				st{"had", "job", 1, "nothing", "you"},
				st{"for", "phore", 4, "four", "fore"},
			},
		}, NewList().Append(NewTestStruct(
			NewBytes([]byte("you")),
			NewBytes([]byte("had")),
			NewNumber(1),
			NewBytes([]byte("job")),
		)).Append(NewTestStruct(
			NewBytes([]byte("fore")),
			NewBytes([]byte("for")),
			NewNumber(4),
			NewBytes([]byte("phore")),
		)), false},
		{"int", args{int(1)}, NewNumber(1), false},
		{"int64", args{int64(1000)}, NewNumber(1000), false},
		{"minus 1", args{int(-1)}, NewNumber(-1), false},
		{"uint64", args{uint64(1)}, NewNumber(1), false},
		{"illegal uint64", args{uint64(math.MaxUint64)}, nil, true},
		{"string", args{"hello"}, NewBytes([]byte("hello")), false},
		{"[]byte", args{[]byte("hello")}, NewBytes([]byte("hello")), false},
		{"time", args{tt}, ts, false},
		{"true", args{true}, NewNumber(1), false},
		{"false", args{false}, NewNumber(0), false},
		{"[]int", args{[]int{1, 23}}, NewList().Append(NewNumber(1)).Append(NewNumber(23)), false},
		{"[] illegal values", args{[]int32{1, 23}}, nil, true},
		{"map", args{map[int]int{1: 2}}, nil, true},
		{"ptr to time", args{&tt}, ts, false},
		{"[][]int", args{[][]int{[]int{1}, []int{2, 3}}},
			NewList().Append(NewList().Append(NewNumber(1))).Append(NewList().Append(NewNumber(2)).Append(NewNumber(3))), false},
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
	type args struct {
		x interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{"simple", args{
			struct {
				X int `chain:"0"`
			}{3},
		}, map[string]int{"X": 0}, false},
		{"in order", args{
			struct {
				X int `chain:"0"`
				Y int `chain:"1"`
				Z int `chain:"2"`
			}{3, 4, 5},
		}, map[string]int{"X": 0, "Y": 1, "Z": 2}, false},
		{"out of order", args{
			struct {
				X int `chain:"2"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, map[string]int{"X": 2, "Y": 0, "Z": 1}, false},
		{"not continuous should not error", args{
			struct {
				X int `chain:"3"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, map[string]int{"X": 3, "Y": 0, "Z": 1}, false},
		{"mixed types", args{
			struct {
				X string `chain:"0"`
				Y int64  `chain:"1"`
				Z byte   `chain:"2"`
			}{"hi", math.MaxInt64, 0x2A},
		}, map[string]int{"X": 0, "Y": 1, "Z": 2}, false},
		{"rename", args{
			struct {
				X int `chain:"0,foo"`
				Y int `chain:"1,Bar"`
				Z int `chain:"2"`
			}{3, 4, 5},
		}, map[string]int{"FOO": 0, "BAR": 1, "Z": 2}, false},
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
		}, map[string]int{"X": 0, "Z": 1}, false},
		{"no chain tags", args{
			struct {
				X int
				Y int
				Z int
			}{3, 4, 5},
		}, nil, true},
		{"not a struct", args{[]int{3}}, nil, true},
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
