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
		}, NewStruct(NewNumber(3)), false},
		{"in order", args{
			struct {
				X int `chain:"0"`
				Y int `chain:"1"`
				Z int `chain:"2"`
			}{3, 4, 5},
		}, NewStruct(NewNumber(3), NewNumber(4), NewNumber(5)), false},
		{"out of order", args{
			struct {
				X int `chain:"2"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, NewStruct(NewNumber(4), NewNumber(5), NewNumber(3)), false},
		{"not continuous should error", args{
			struct {
				X int `chain:"3"`
				Y int `chain:"0"`
				Z int `chain:"1"`
			}{3, 4, 5},
		}, nil, true},
		{"mixed types", args{
			struct {
				X string    `chain:"0"`
				Y int64     `chain:"1"`
				Z byte      `chain:"2"`
				T time.Time `chain:"3"`
			}{"hi", math.MaxInt64, 0x2A, tt},
		}, NewStruct(NewBytes([]byte("hi")), NewNumber(math.MaxInt64), NewNumber(42), ts), false},
		{"simple array", args{
			[]int{1, 2, 3},
		}, NewList().Append(NewNumber(1)).Append(NewNumber(2)).Append(NewNumber(3)), false},
		{"struct array", args{
			[]st{
				st{"had", "job", 1, "nothing", "you"},
				st{"for", "phore", 4, "four", "fore"},
			},
		}, NewList().Append(NewStruct(
			NewBytes([]byte("you")),
			NewBytes([]byte("had")),
			NewNumber(1),
			NewBytes([]byte("job")),
		)).Append(NewStruct(
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
		{"map", args{map[int]int{1: 2}}, nil, true},
		{"ptr to time", args{&tt}, ts, false},
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
