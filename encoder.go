package borsh

import (
	"io"
	"reflect"
)

type Encoder struct {
	w io.Writer
	p *pool
}

func NewEncoder(w io.Writer) *Encoder {
	p := newPool()
	return &Encoder{w: w, p: p}
}

func (e *Encoder) Encode(s interface{}) error {
	return serialize(reflect.ValueOf(s), e.w, e.p)
}

func (e *Encoder) Close() error {
	return nil
}
