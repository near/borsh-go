# borsh-go

## Example

```go
import "github.com/ouromoros/borsh-go"

type A struct {
    x uint64,
    y string,
}

func testSimple() {
    x := A{
        x: 3301,
        y: "liber primus",
    }
	data, err := borsh.Serialize(x)
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
