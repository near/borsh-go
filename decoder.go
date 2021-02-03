package borsh

import (
	"errors"
	"io"
	"reflect"
)

type Decoder struct {
	r io.Reader
	p *pool
}

func NewDecoder(r io.Reader) *Decoder {
	p := newPool()
	return &Decoder{r: r, p: p}
}

func (d *Decoder) Decode(s interface{}) error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return errors.New("argument must be pointer")
	}
	val, err := deserialize(t, d.r, d.p)
	if err != nil {
		return nil
	}
	reflect.ValueOf(s).Elem().Set(reflect.ValueOf(val))
	return nil
}

func (d *Decoder) Close() error {
	return nil
}
