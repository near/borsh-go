package borsh

import (
	"bytes"
	"math"
	"math/big"
	"reflect"
	strings2 "strings"
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

type C struct {
	A3 [3]int
	S  []int
	P  *int
	M  map[string]string
}

func TestBasicContainer(t *testing.T) {
	ip := new(int)
	*ip = 213
	x := C{
		A3: [3]int{234, -123, 123},
		S:  []int{21442, 421241241, 2424},
		P:  ip,
		M:  make(map[string]string),
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

type N struct {
	B B
	C C
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
			M:  make(map[string]string),
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

func TestSimpleEnum(t *testing.T) {
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

type ComplexEnum struct {
	Enum Enum `borsh_enum:"true"`
	Foo  Foo
	Bar  Bar
}

type Foo struct {
	FooA int32
	FooB string
}

type Bar struct {
	BarA int64
	BarB string
}

func TestComplexEnum(t *testing.T) {
	x := ComplexEnum{
		Enum: 0,
		Foo: Foo{
			FooA: 23,
			FooB: "baz",
		},
	}
	data, err := Serialize(x)
	if err != nil {
		t.Fatal(err)
	}
	y := new(ComplexEnum)
	err = Deserialize(y, data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Fatal(x, y)
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

type Skipped struct {
	A int
	B int `borsh_skip:"true"`
	C int
}

func TestSkipped(t *testing.T) {
	x := Skipped{
		A: 32,
		B: 535,
		C: 123,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(Skipped)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if x.A != y.A || x.C != y.C {
		t.Errorf("%v fields not equal to %v", x, y)
	}
	if y.B == x.B {
		t.Errorf("didn't skip field B")
	}
}

type E struct{}

func TestEmpty(t *testing.T) {
	x := E{}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 0 {
		t.Error("not empty")
	}
	y := new(E)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

func testValue(t *testing.T, v interface{}) {
	data, err := Serialize(v)
	if err != nil {
		t.Error(err)
	}
	parsed := reflect.New(reflect.TypeOf(v))
	err = Deserialize(parsed.Interface(), data)
	if err != nil {
		t.Error(err)
	}
	reflect.DeepEqual(v, parsed.Elem().Interface())
}

func TestStrings(t *testing.T) {
	tests := []struct {
		in string
	}{
		{""},
		{"a"},
		{"hellow world"},
		{strings2.Repeat("x", 1024)},
		{strings2.Repeat("x", 4096)},
		{strings2.Repeat("x", 65535)},
		{strings2.Repeat("hello world!", 1000)},
		{"💩"},
	}

	for _, tt := range tests {
		testValue(t, tt.in)
	}
}

func makeInt32Slice(val int32, len int) []int32 {
	s := make([]int32, len)
	for i := 0; i < len; i++ {
		s[i] = val
	}
	return s
}

func TestSlices(t *testing.T) {
	tests := []struct {
		in []int32
	}{
		{makeInt32Slice(1000000000, 0)},
		{makeInt32Slice(1000000000, 1)},
		{makeInt32Slice(1000000000, 2)},
		{makeInt32Slice(1000000000, 3)},
		{makeInt32Slice(1000000000, 4)},
		{makeInt32Slice(1000000000, 8)},
		{makeInt32Slice(1000000000, 16)},
		{makeInt32Slice(1000000000, 32)},
		{makeInt32Slice(1000000000, 64)},
		{makeInt32Slice(1000000000, 65)},
	}

	for _, tt := range tests {
		testValue(t, tt.in)
	}
}

func TestUint128(t *testing.T) {
	tests := []struct {
		in big.Int
	}{
		{*big.NewInt(23)},
		{*big.NewInt(math.MaxInt16)},
		{*big.NewInt(math.MaxInt32)},
		{*big.NewInt(math.MaxInt64)},
		{*big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))},
	}

	for _, tt := range tests {
		testValue(t, tt.in)
	}
}

type Myu8 uint8
type Myu16 uint16
type Myu32 uint32
type Myu64 uint64
type Myi8 int8
type Myi16 int16
type Myi32 int32
type Myi64 int64

type CustomType struct {
	U8  Myu8
	U16 Myu16
	U32 Myu32
	U64 Myu64
	I8  Myi8
	I16 Myi16
	I32 Myi32
	I64 Myi64
}

func TestCustomType(t *testing.T) {
	x := CustomType{
		U8:  1,
		U16: 2,
		U32: 3,
		U64: 4,
		I8:  5,
		I16: 6,
		I32: 7,
		I64: 8,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(CustomType)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type BoolStruct struct {
	T bool
	F bool
}

func TestBool(t *testing.T) {
	x := BoolStruct{
		T: true,
		F: false,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(BoolStruct)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

func TestMap(t *testing.T) {
	type Key [32]uint8
	m := make(map[Key]uint8)
	key1 := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134, 9, 1}
	key2 := Key{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134, 9, 1}
	key3 := Key{1, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134}
	m[key1] = 1
	m[key2] = 2
	m[key3] = 3
	bts, err := Serialize(m)
	if err != nil {
		t.Error(err)
	}
	n := make(map[Key]uint8)
	err = Deserialize(&n, bts)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(m, n) {
		t.Error(m, n)
	}

}

func TestInterface(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want []byte
	}{
		{
			name: "serialize struct",
			data: struct {
				Data uint8
			}{
				Data: 1,
			},
			want: []byte{1},
		},
		{
			name: "serialize uint8 array",
			data: struct {
				Data []uint8
			}{
				Data: []uint8{1, 2, 3, 4},
			},
			want: []byte{4, 0, 0, 0, 1, 2, 3, 4},
		},
		{
			name: "serialize interface array",
			data: struct {
				Data []interface{}
			}{
				Data: []interface{}{
					uint8(1),
					uint16(2),
					uint32(3),
					uint64(4),
					struct {
						Data uint8
					}{
						Data: 5,
					},
					struct {
						Data []uint8
					}{
						Data: []uint8{6, 7, 8, 9},
					},
				},
			},
			want: []byte{6, 0, 0, 0, 1, 2, 0, 3, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 5, 4, 0, 0, 0, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := Serialize(tt.data)
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(b, tt.want) {
				t.Errorf("want: %v, got: %v", tt.want, b)
			}
		})
	}
}
