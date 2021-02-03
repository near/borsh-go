package borsh

type pool struct {
	b1 []byte
	b2 []byte
	b4 []byte
	b8 []byte
}

func newPool() *pool {
	p := new(pool)
	p.b1 = make([]byte, 1)
	p.b2 = make([]byte, 2)
	p.b4 = make([]byte, 4)
	p.b8 = make([]byte, 8)
	return p
}

func (p *pool) getBytes(n int) []byte {
	switch n {
	case 1:
		return p.b1
	case 2:
		return p.b2
	case 4:
		return p.b4
	case 8:
		return p.b8
	}
	panic("unsupported number of bytes")
}
