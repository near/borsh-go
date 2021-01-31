package borsh

import (
	"reflect"
	"testing"
)

type A struct {
	A int
	B int32
}

func TestSimple(t *testing.T) {
	x := A{
		A: 1,
		B: 32,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(A)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}
