package borsh_test

import (
	"github.com/near/borsh-go"
	"testing"
)

func TestMap(t *testing.T) {
	type Key [32]uint8
	m := make(map[Key]uint8)
	key1 := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134, 9, 1}
	key2 := Key{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134, 9, 1}
	key3 := Key{1, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134}
	m[key1] = 1
	m[key2] = 2
	m[key3] = 3
	bytes, err := borsh.Serialize(m)
	if err != nil {
		t.Fatal(err)
	}

	m2 := make(map[Key]uint8)
	err = borsh.Deserialize(&m2, bytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m)
}

func TestMap2(t *testing.T) {
	type Key [32]uint8
	m := make(map[Key]uint8)
	key1 := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134, 9, 1}
	// key2 := Key{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134, 9, 1}
	// key3 := Key{1, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 169, 224, 9, 91, 137, 101, 192, 30, 106, 9, 201, 121, 56, 243, 134}
	m[key1] = 1
	// m[key2] = 1
	// m[key3] = 1
	borsh.Serialize(m)
}
