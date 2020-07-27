/*
"Missing feature" of the Go packages `encoding` and `fmt`: unmarshal arbitrary
text into an arbitrary value. Counterpart to the marshaling functionality of
`fmt.Sprintf("%v")`.

Examples

Decode individual values:

	var num int64
	err := untext.UnmarshalString(`10`, &num)

	var inst time.Time
	err = untext.UnmarshalString(`0001-02-03T04:05:06Z`, &inst)

Decode slices:

	var nums []int64
	err = untext.UnmarshalSlice([]string{`10`, `20`}, &nums)
*/
package untext

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

/*
Unmarshals arbitrary text into an arbitrary destination pointer. Supports a
variety of "well-known" types out of the box, and falls back on
`encoding.TextUnmarshaler`.
*/
func UnmarshalBytes(input []byte, dest interface{}) error {
	impl, _ := dest.(encoding.TextUnmarshaler)
	if impl != nil {
		err := impl.UnmarshalText(input)
		return err
	}

	rval, err := settableRval(dest)
	if err != nil {
		return err
	}

	return unmarshalRval(input, rval)
}

/*
Variant of `UnmarshalBytes` that accepts a string as input.
*/
func UnmarshalString(input string, dest interface{}) error {
	return UnmarshalBytes(stringToBytesUnsafe(input), dest)
}

var unmarshalerRtype = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func unmarshalRval(input []byte, rval reflect.Value) error {
	if !rval.CanSet() {
		return fmt.Errorf(`expected settable rval, got %v`, rval.Interface())
	}
	rtype := rval.Type()

	switch rval.Type() {
	// Time is a special case for symmetry with JSON and SQL.
	case timeRtype:
		var inst time.Time
		err := inst.UnmarshalText(input)
		if err != nil {
			return err
		}
		rval.Set(reflect.ValueOf(inst))
		return nil
	}

	inputStr := bytesToMutableString(input)

	switch rtype.Kind() {
	case reflect.Bool:
		switch inputStr {
		case "true":
			rval.SetBool(true)
			return nil
		case "false":
			rval.SetBool(false)
			return nil
		default:
			return fmt.Errorf(`can't unmarshal %q as bool`, input)
		}

	case reflect.Int8:
		num, err := strconv.ParseInt(inputStr, 10, 8)
		rval.SetInt(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Int16:
		num, err := strconv.ParseInt(inputStr, 10, 16)
		rval.SetInt(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Int32:
		num, err := strconv.ParseInt(inputStr, 10, 32)
		rval.SetInt(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Int64:
		num, err := strconv.ParseInt(inputStr, 10, 64)
		rval.SetInt(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Uint8:
		num, err := strconv.ParseUint(inputStr, 10, 8)
		rval.SetUint(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Uint16:
		num, err := strconv.ParseUint(inputStr, 10, 16)
		rval.SetUint(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Uint32:
		num, err := strconv.ParseUint(inputStr, 10, 32)
		rval.SetUint(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Uint64:
		num, err := strconv.ParseUint(inputStr, 10, 64)
		rval.SetUint(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Float32:
		num, err := strconv.ParseFloat(inputStr, 32)
		rval.SetFloat(num)
		return maybeUnmarshalErr(input, err)

	case reflect.Float64:
		num, err := strconv.ParseFloat(inputStr, 64)
		rval.SetFloat(num)
		return maybeUnmarshalErr(input, err)

	case reflect.String:
		rval.SetString(bytesToStringAlloc(input))
		return nil

	case reflect.Ptr:
		ptrRval := reflect.New(rval.Type().Elem())

		if ptrRval.Type().Implements(unmarshalerRtype) {
			impl := ptrRval.Interface().(encoding.TextUnmarshaler)
			err := impl.UnmarshalText(input)
			if err != nil {
				return err
			}
			rval.Set(ptrRval)
			return nil
		}

		err := unmarshalRval(input, ptrRval.Elem())
		if err != nil {
			return err
		}

		rval.Set(ptrRval)
		return nil

	// Missing:
	// case reflect.Array
	// case reflect.Slice
	// case reflect.Struct

	default:
		return fmt.Errorf(`can't unmarshal %q into type %q`, input, rtype)
	}
}

func maybeUnmarshalErr(input []byte, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(`failed to unmarshal %q: %+v`, input, err)
}

var timeRtype = reflect.TypeOf(time.Time{})

/*
Unmarshals a slice of strings by calling `UnmarshalString` for each value. The
destination must be a non-nil pointer to a slice. The element type determines
the output type for `UnmarshalString`.
*/
func UnmarshalSlice(inputs []string, dest interface{}) error {
	rval, err := settableSliceRval(dest)
	if err != nil {
		return err
	}

	rval.Set(reflect.MakeSlice(rval.Type(), len(inputs), len(inputs)))

	for i, input := range inputs {
		err := UnmarshalString(input, rval.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func settableRval(input interface{}) (reflect.Value, error) {
	rval := reflect.ValueOf(input)
	rtype := rval.Type()
	if rtype.Kind() != reflect.Ptr {
		return rval, fmt.Errorf(`expected a pointer, got a %q`, rtype)
	}
	rval = rval.Elem()
	if !rval.CanSet() {
		return rval, fmt.Errorf(`can't set into non-settable value of type %q`, rtype)
	}
	return rval, nil
}

func settableSliceRval(dest interface{}) (reflect.Value, error) {
	rval, err := settableRval(dest)
	if err != nil {
		return rval, err
	}
	if rval.Type().Kind() != reflect.Slice {
		return rval, fmt.Errorf(`expected a slice pointer, got a %q`, rval.Type())
	}
	return rval, nil
}

// Self-reminder about non-free conversions.
func bytesToStringAlloc(bytes []byte) string { return string(bytes) }

// Self-reminder about non-free conversions.
func stringToBytesAlloc(input string) []byte { return []byte(input) }

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Reasonably safe. Should not be used when the underlying
byte array is volatile, for example when it's part of a scratch buffer during
SQL scanning.
*/
func bytesToMutableString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

/*
Allocation-free conversion. Returns a byte slice backed by the provided string.
Mutations are reflected in the source string, unless it's backed by constant
storage, in which case they trigger a segfault. Reslicing is ok. Should be safe
as long as the resulting bytes are not mutated. Sometimes produces unexpected
garbage, possibly because the string was, in turn, backed by mutable storage
which gets modified before we use the result; needs investigation.
*/
func stringToBytesUnsafe(input string) []byte {
	type sliceHeader struct {
		dat uintptr
		len int
		cap int
	}
	slice := *(*sliceHeader)(unsafe.Pointer(&input))
	slice.cap = slice.len
	return *(*[]byte)(unsafe.Pointer(&slice))
}
