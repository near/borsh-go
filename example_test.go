package borsh_test

import "github.com/ouromoros/borsh-go"

type A struct {
	B int32
}

func ExampleDeserialize() {
	a := A{B: 123}
	data, _ := borsh.Serialize(a)
	b := &A{}
	borsh.Deserialize(b, data)
}

func ExampleSerialize() {
	a := A{B: 123}
	data, _ := borsh.Serialize(a)
	print(data)
}
