# borsh-go

**borsh-go** is an implementation of the [Borsh] binary serialization format for Go
projects.

Borsh stands for _Binary Object Representation Serializer for Hashing_. It is
meant to be used in security-critical projects as it prioritizes consistency,
safety, speed, and comes with a strict specification.

## Features

- Based on Go Reflection. Avoids the need for create protocol file and code generation. Simply
defining `struct` and go.


## Type Mappings

Borsh                 | Go           |  Description
--------------------- | -------------- |--------
`u7` integer          | `uint8`        | 
`u15` integer         | `uint16`       |
`u31` integer         | `uint32`       |
`u63` integer         | `uint64`       |
`u127` integer        |            |  Not supported yet
`i7` integer          | `int8`        | 
`i15` integer         | `int16`       |
`i31` integer         | `int32`       |
`i63` integer         | `int64`       |
`i127` integer        |            |  Not supported yet
`f31` float           | `float32`      |
`f63` float           | `float64`      |
fixed-size array      | `[size]type`   |  go array
dynamic-size array    |  `[]type`      |  go slice
string                | `string`       |
option                |  `*type`         |   go pointer
map                   |   `map`          |
set                   |   `map[type]struct{}`  | go map with value type set to `struct{}`
structs               |   `struct`      |
enum                  |   `borsh.Enum`  |    use `type MyEnum borsh.Enum` to define


## Usage

### Example

```go
package demo

import (
	"github.com/ouromoros/borsh-go"
	"log"
	"reflect"
	"testing"
)

type A struct {
	X uint64
	Y string
    Z string `borsh_skip:"true"` // will skip this field when serializing/deserializing
}

func TestSimple(t *testing.T) {
	x := A{
		X: 3301,
		Y: "liber primus",
	}
	data, err := borsh.Serialize(x)
	log.Print(data)
	if err != nil {
		t.Error(err)
	}
	y := new(A)
	err = borsh.Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}
```
