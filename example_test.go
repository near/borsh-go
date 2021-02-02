package borsh_test

import (
	"bytes"
	"github.com/ouromoros/borsh-go"
	"log"
	"strings"
)

type A struct {
	B int32
}

func ExampleDeserialize() {
	a := A{B: 123}
	data, _ := borsh.Serialize(a)
	b := &A{}
	err := borsh.Deserialize(b, data)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSerialize() {
	a := A{B: 123}
	data, err := borsh.Serialize(a)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(data)
}

func ExampleNewEncoder() {
	a := A{B: 123}
	b := strings.Builder{}
	e := borsh.NewEncoder(&b)
	if err := e.Encode(a); err != nil {
		log.Fatal(err)
	}
}

func ExampleNewDecoder() {
	a := A{B: 123}
	data, _ := borsh.Serialize(a)
	b := &A{}
	d := borsh.NewDecoder(bytes.NewReader(data))
	if err := d.Decode(b); err != nil {
		log.Fatal(err)
	}
}
