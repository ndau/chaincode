package converter

import (
	"errors"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/oneiro-ndev/signature/pkg/signature"
)

// chain tags have comma-separated values. The first is always a numeric field index, the second
// is an optional field name to be used in CHASM (if the second is not specified, the actual
// field name is used). Names are converted to uppercase for use in the assembler. Names, when
// converted to uppercase, must be valid CHASM constant names ([A-Z][A-Z0-9_]*)

// parseChainTag interprets a tag string.
func parseChainTag(tag string, name string) (byte, string, error) {
	sp := strings.Split(tag, ",")
	ix, err := strconv.ParseInt(sp[0], 10, 8)
	if err != nil {
		return 0, "", err
	}
	if len(sp) > 1 {
		p := regexp.MustCompile("[A-Za-z][A-Za-z0-9_]+")
		if !p.MatchString(sp[1]) {
			return 0, "", errors.New("name must be a valid constant in chasm ([A-Z][A-Z0-9_]*)")
		}
		name = sp[1]
	}
	return byte(ix), strings.ToUpper(name), nil
}

// ToValueScalar converts a scalar value to a VM Value object
// it handles ints of several types, bool, string, time.Time, and
// pointers to these.
func ToValueScalar(x interface{}) (vm.Value, error) {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return vm.NewTrue(), nil
		}
		return vm.NewFalse(), nil
	case reflect.Int, reflect.Int64:
		n := v.Int()
		return vm.NewNumber(n), nil
	case reflect.Uint64, reflect.Uint8:
		n := v.Uint()
		if n > math.MaxInt64 {
			return nil, errors.New("value doesn't fit into a Number")
		}
		return vm.NewNumber(int64(n)), nil
	case reflect.String:
		return vm.NewBytes([]byte(x.(string))), nil
	case reflect.Ptr:
		// convert pointers to the object they point to and try again recursively
		return ToValueScalar(v.Elem().Interface())
	case reflect.Struct:
		// if we get a struct at this level, we have to see if it is one of our
		// special types we already understand
		switch v.Type() {
		case reflect.ValueOf(time.Time{}).Type():
			return vm.NewTimestampFromTime(x.(time.Time))
		case reflect.ValueOf(address.Address{}).Type():
			data, err := x.(address.Address).MarshalText()
			if err != nil {
				return nil, err
			}
			return vm.NewBytes(data), nil
		case reflect.ValueOf(signature.Signature{}).Type():
			data, err := x.(signature.Signature).Marshal()
			if err != nil {
				return nil, err
			}
			return vm.NewBytes(data), nil
		default:
			// try calling the struct parser recursively
			level2, err := ToValue(x)
			if err != nil {
				return nil, err
			}
			return level2, nil
		}
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
// Structs are treated recursively; all field IDs must be distinct at all levels, because
// the generated struct is flat.
// Arrays create a list of values in the array
func ToValue(x interface{}) (vm.Value, error) {
	vx := reflect.ValueOf(x)
	tx := reflect.TypeOf(x)
	switch vx.Kind() {
	case reflect.Array, reflect.Slice:
		// special case for byte arrays -- they are treated as strings
		if tx == reflect.TypeOf([]byte{}) {
			return ToValueScalar(string(x.([]byte)))
		}
		// if it's an array, create a list out of the individual items by calling this function
		// recursively. This will also work for arrays of arrays or arrays of structs.
		li := vm.NewList()
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
			return vm.NewTimestampFromTime(x.(time.Time))
		}

		// if it's a struct, iterate the members and look to see if they have "chain:" tags;
		// if so, assemble a struct from all the members that do. If no chain tags exist, then
		// error.
		st := vm.NewStruct()
		for i := 0; i < tx.NumField(); i++ {
			fld := tx.Field(i)
			tag := fld.Tag.Get("chain")

			ix, _, err := parseChainTag(tag, "")
			if err != nil {
				return nil, err
			}

			child, err := ToValueScalar(vx.FieldByIndex(fld.Index).Interface())
			if err != nil {
				return nil, err
			}
			if chstr, ok := child.(*vm.Struct); ok {
				for _, ind := range chstr.Indices() {
					v, _ := chstr.Get(ind)
					st, err = st.SafeSet(ind, v)
					if err != nil {
						return nil, err
					}
				}
				continue
			}

			// if there's no chain tag (and it wasn't a struct), just move on
			if ix < 0 {
				continue
			}

			st, err = st.SafeSet(ix, child)
			if err != nil {
				return nil, err
			}
		}
		return st, nil

	default:
		// for all other types assume it's a scalar
		return ToValueScalar(x)
	}
}

// ExtractConstants takes an interface which should be a Go language struct with
// "chain" Struct Tags, and extracts a map of names to indices in the generated vm struct
func ExtractConstants(x interface{}) (map[string]byte, error) {
	vx := reflect.ValueOf(x)
	tx := reflect.TypeOf(x)
	switch vx.Kind() {
	case reflect.Struct:
		// if it's a struct, iterate the members and look to see if they have "chain:" tags;
		// if so, assemble a map from all the members that do. If no chain tags exist, then
		// error.
		result := make(map[string]byte)
		for i := 0; i < tx.NumField(); i++ {
			fld := tx.Field(i)
			tag := fld.Tag.Get("chain")
			// if there's no chain tag, just move on
			if tag == "" {
				continue
			}

			ix, name, err := parseChainTag(tag, fld.Name)
			if err != nil {
				return nil, err
			}

			result[name] = ix
		}
		if len(result) == 0 {
			return nil, errors.New("no chain tags found in struct")
		}
		return result, nil

	default:
		// all other types are an error
		return nil, errors.New("object was not a tagged struct")
	}
}
