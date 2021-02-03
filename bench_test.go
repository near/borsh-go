package borsh_test

import (
	"github.com/ouromoros/borsh-go"
	"math/rand"
	"testing"
)

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

func BenchmarkDeserialize(t *testing.B) {
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
	data, err := borsh.Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(B)
	for i := 0; i < t.N; i++ {
		err = borsh.Deserialize(y, data)
		if err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkSerialize(t *testing.B) {
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
	for i := 0; i < t.N; i++ {
		_, _ = borsh.Serialize(x)
	}
}

func BenchmarkFuzzSerialize(t *testing.B) {
	s1 := rand.NewSource(42)
	r1 := rand.New(s1)

	for i := 0; i < 100; i++ {
		st := fuzzType(r1, 0)
		val := fuzzValue(r1, st)
		for j := 0; j < t.N; j++ {
			_, _ = borsh.Serialize(val)
		}
	}
}

func BenchmarkFuzzDeserialize(t *testing.B) {
	s1 := rand.NewSource(42)
	r1 := rand.New(s1)

	for i := 0; i < 100; i++ {
		st := fuzzType(r1, 0)
		val := fuzzValue(r1, st)
		data, _ := borsh.Serialize(val)
		for j := 0; j < t.N; j++ {
			_ = borsh.Deserialize(val, data)
		}
	}
}
