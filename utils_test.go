package photoshare

import (
	"testing"
)

func TestPgArrToIntSlice(t *testing.T) {
	a := "{1,2,3,4,5}"
	result := pgArrToIntSlice(a)
	if len(result) != 5 {
		t.Fail()
	}
}

func TestIntSliceToPgArr(t *testing.T) {
	s := []int64{1, 2, 3, 4, 5}
	result := intSliceToPgArr(s)
	if result != "{1,2,3,4,5}" {
		t.Fail()
	}
}
