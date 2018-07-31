package vm

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// chain tags have comma-separated values. The first is always a numeric field index, the second
// is an optional field name to be used in CHASM (if the second is not specified, the actual
// field name is used). Names are converted to uppercase for use in the assembler.

// parseChainTag interprets a tag string
func parseChainTag(tag string, name string) (int, string, error) {
	if tag == "" {
		return -1, "", nil
	}
	sp := strings.Split(tag, ",")
	ix, err := strconv.ParseInt(sp[0], 10, 8)
	if err != nil {
		return 0, "", err
	}
	if len(sp) > 1 {
		name = sp[1]
	}
	return int(ix), strings.ToUpper(name), nil
}

// ToValueScalar converts a scalar value to a VM Value object
// it handles ints of several types, bool, string, time.Time, and
// pointers to these.
func ToValueScalar(x interface{}) (Value, error) {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Bool:
		if x.(bool) {
			return NewNumber(1), nil
		}
		return NewNumber(0), nil
	case reflect.Int:
		n := int64(x.(int))
		return NewNumber(n), nil
	case reflect.Int64:
		n := x.(int64)
		return NewNumber(n), nil
	case reflect.Uint64:
		n := x.(uint64)
		if n > math.MaxInt64 {
			return nil, errors.New("value doesn't fit into a Number")
		}
		return NewNumber(int64(n)), nil
	case reflect.Uint8:
		n := int64(x.(uint8))
		return NewNumber(n), nil
	case reflect.String:
		return NewBytes([]byte(x.(string))), nil
	case reflect.Ptr:
		// convert pointers to the object they point to and try again recursively
		return ToValueScalar(v.Elem().Interface())
	case reflect.Struct:
		// if we get a struct at this level it should be a timestamp
		if v.Type() == reflect.ValueOf(time.Time{}).Type() {
			// it's a time.Time so convert it to a Timestamp
			return NewTimestampFromTime(x.(time.Time))
		}
		return nil, errors.New("is struct, not a scalar")
	case reflect.Interface:
		// we can't handle generic interfaces
		return nil, errors.New("is interface, not a scalar")
	case reflect.Array, reflect.Map, reflect.Slice:
		// and arrays and slices happen at a higher level
		return nil, errors.New("is container, not a scalar")
	}
	return nil, errors.New("unknown type")
}

// ToValue returns a Go value as a VM value, including if the Go value is a struct or array.
// Structs are not treated recursively; only the top level is examined.
// Arrays create a list of values in the array
func ToValue(x interface{}) (Value, error) {
	vx := reflect.ValueOf(x)
	tx := reflect.TypeOf(x)
	switch vx.Kind() {
	case reflect.Array, reflect.Slice:
		if tx == reflect.TypeOf([]byte{}) {
			return ToValueScalar(string(x.([]byte)))
		}
		// if it's an array, create a list out of the individual items by calling this function
		// recursively. This will work for arrays of arrays or arrays of structs.
		li := NewList()
		for i := 0; i < vx.Len(); i++ {
			item := vx.Index(i).Interface()
			v, err := ToValue(item)
			if err != nil {
				return nil, err
			}
			li = li.Append(v)
		}
		return li, nil

	case reflect.Struct:
		// first check to see if it's just a timestamp
		if vx.Type() == reflect.ValueOf(time.Time{}).Type() {
			// it's a time.Time so convert it to a Timestamp
			return NewTimestampFromTime(x.(time.Time))
		}

		// if it's a struct, iterate the members and look to see if they have "chain:" tags;
		// if so, assemble a struct from all the members that do. If no chain tags exist, then
		// error.
		fm := make(map[int]Value)
		for i := 0; i < tx.NumField(); i++ {
			fld := tx.Field(i)
			tag := fld.Tag.Get("chain")

			ix, _, err := parseChainTag(tag, "")
			// if there's no chain tag, just move on
			if ix < 0 {
				continue
			}
			if err != nil {
				return nil, err
			}

			fm[ix], err = ToValueScalar(vx.FieldByIndex(fld.Index).Interface())
			if err != nil {
				return nil, err
			}
		}
		// if we get here, we have a map of indices and values; we want that map to have
		// indices from 0 to len(fm)-1 or the tags are messed up
		st := NewStruct()
		for i := 0; i < len(fm); i++ {
			v, ok := fm[i]
			if !ok {
				return nil, errors.New("struct indices were not adjacent")
			}
			st = st.Append(v)
		}
		return st, nil

	default:
		// for all other types assume it's a scalar
		return ToValueScalar(x)
	}
}
