package borsh

import (
	"reflect"
	"testing"
)

type A struct {
	A int
	B int32
}

type B struct {
	I8  int8
	I16 int16
	I32 int32
	I64 int64
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
	F32 float32
	F64 float64
}

type C struct {
	A3 [3]int
	S  []int
	P  *int
}

type N struct {
	B B
	C C
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

func TestBasic(t *testing.T) {
	x := B{
		I8:  12,
		I16: -1,
		I32: 124,
		I64: 1243,
		U8:  1,
		U16: 979,
		U32: 123124,
		U64: 1135351135,
		F32: -231.23,
		F64: 3121221.232,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(B)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

func TestBasicContainer(t *testing.T) {
	ip := new(int)
	*ip = 213
	x := C{
		A3: [3]int{234, -123, 123},
		S:  []int{21442, 421241241, 2424},
		P:  ip,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(C)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

func TestNested(t *testing.T) {
	ip := new(int)
	*ip = 213
	x := N{
		B: B{
			I8:  12,
			I16: -1,
			I32: 124,
			I64: 1243,
			U8:  1,
			U16: 979,
			U32: 123124,
			U64: 1135351135,
			F32: -231.23,
			F64: 3121221.232,
		},
		C: C{
			A3: [3]int{234, -123, 123},
			S:  []int{21442, 421241241, 2424},
			P:  ip,
		},
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(N)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type Dummy Enum

const (
	x Dummy = iota
	y
	z
)

type D struct {
	D Dummy
}

func TestEnum(t *testing.T) {
	x := D{
		D: y,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(D)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type S struct {
	S map[int]struct{}
}

func TestSet(t *testing.T) {
	x := S{
		S: map[int]struct{}{124: struct{}{}, 214: struct{}{}, 24: struct{}{}, 53: struct{}{}},
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(S)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}
