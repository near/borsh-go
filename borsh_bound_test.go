package borsh

import (
	"runtime"
	"testing"
)

// A tiny payload that declares a ~512 MiB string but supplies no bytes for it.
// Before the fix, Deserialize allocated the full declared length up front (via
// make([]byte, l)) and only then noticed the data was truncated — a trivial
// out-of-memory DoS from a 4-byte input. After the fix it must reject in O(1)
// without a large allocation.
func TestDeserializeStringRejectsLengthBeyondRemaining(t *testing.T) {
	type S struct{ Name string }
	data := []byte{0x00, 0x00, 0x00, 0x20} // LE uint32 = 0x20000000 = 512 MiB

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	var s S
	err := Deserialize(&s, data)

	runtime.ReadMemStats(&after)
	alloc := int64(after.TotalAlloc - before.TotalAlloc)
	t.Logf("err=%v  bytes allocated during call=%d", err, alloc)

	if err == nil {
		t.Fatalf("expected an error for a truncated string, got nil")
	}
	if alloc > 10*1024*1024 {
		t.Fatalf("over-allocated %d bytes before validating the declared length (DoS)", alloc)
	}
}

// Honestly-encoded input must still round-trip unchanged after the fix.
func TestDeserializeStillAcceptsValidInput(t *testing.T) {
	type S struct {
		Name  string
		Nums  []uint32
		Attrs map[string]uint8
	}
	in := S{
		Name:  "near",
		Nums:  []uint32{1, 2, 3},
		Attrs: map[string]uint8{"a": 1, "b": 2},
	}
	enc, err := Serialize(in)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	var out S
	if err := Deserialize(&out, enc); err != nil {
		t.Fatalf("Deserialize of valid input failed: %v", err)
	}
	if out.Name != in.Name || len(out.Nums) != 3 || len(out.Attrs) != 2 {
		t.Fatalf("round-trip mismatch: got %+v want %+v", out, in)
	}
}
