# borsh-go

## Example

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
