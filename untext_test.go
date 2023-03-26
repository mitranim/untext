package untext

import (
	"reflect"
	"testing"
	"time"
)

/*
TODO test more types.
*/
func TestParse(t *testing.T) {
	{
		var expected int64 = 10
		var result int64

		err := Parse(`10`, &result)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if expected != result {
			if err != nil {
				t.Fatalf("expected %v, got %v", expected, result)
			}
		}
	}

	{
		var expected time.Time = time.Date(1, 2, 3, 4, 5, 6, 0, time.UTC)
		var result time.Time

		err := Parse(`0001-02-03T04:05:06Z`, &result)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if expected != result {
			if err != nil {
				t.Fatalf("expected %v, got %v", expected, result)
			}
		}
	}
}

/*
TODO test more types.
*/
func TestParseSlice(t *testing.T) {
	{
		var expected = []int64{10, 20}
		var result []int64

		err := ParseSlice([]string{`10`, `20`}, &result)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if !reflect.DeepEqual(expected, result) {
			if err != nil {
				t.Fatalf("expected %v, got %v", expected, result)
			}
		}
	}
}
