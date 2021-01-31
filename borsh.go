package borsh

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"sort"
)

func Deserialize(s interface{}, data []byte) error {
	reader := bytes.NewReader(data)
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("passed struct must be pointer")
	}
	result, err := deserialize(reflect.TypeOf(s).Elem(), reader)
	if err != nil {
		return err
	}
	v.Elem().Set(reflect.ValueOf(result))
	return nil
}

func read(r *bytes.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	l, err := r.Read(b)
	if l != n {
		return nil, errors.New("failed to read required bytes")
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

func deserialize(t reflect.Type, r *bytes.Reader) (interface{}, error) {
	if t.Kind() == reflect.Uint8 {
		tmp, err := read(r, 1)
		if err != nil {
			return nil, err
		}
		e := reflect.New(t)
		e.Elem().Set(reflect.ValueOf(uint8(tmp[0])).Convert(t))
		return e.Elem().Interface(), nil
	}

	switch t.Kind() {
	case reflect.Int8:
		tmp, err := read(r, 1)
		if err != nil {
			return nil, err
		}
		return int8(tmp[0]), nil
	case reflect.Int16:
		tmp, err := read(r, 2)
		if err != nil {
			return nil, err
		}
		return int16(binary.LittleEndian.Uint16(tmp)), nil
	case reflect.Int32:
		tmp, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		return int32(binary.LittleEndian.Uint32(tmp)), nil
	case reflect.Int64:
		tmp, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		return int64(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Int:
		tmp, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		return int(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Uint8:
		tmp, err := read(r, 1)
		if err != nil {
			return nil, err
		}
		return uint8(tmp[0]), nil
	case reflect.Uint16:
		tmp, err := read(r, 2)
		if err != nil {
			return nil, err
		}
		return uint16(binary.LittleEndian.Uint16(tmp)), nil
	case reflect.Uint32:
		tmp, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		return uint32(binary.LittleEndian.Uint32(tmp)), nil
	case reflect.Uint64:
		tmp, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		return uint64(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Uint:
		tmp, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		return uint(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Float32:
		tmp, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		bits := binary.LittleEndian.Uint32(tmp)
		f := math.Float32frombits(bits)
		if math.IsNaN(float64(f)) {
			return nil, errors.New("NaN for float not allowed")
		}
		return f, nil
	case reflect.Float64:
		tmp, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		bits := binary.LittleEndian.Uint64(tmp)
		f := math.Float64frombits(bits)
		if math.IsNaN(f) {
			return nil, errors.New("NaN for float not allowed")
		}
		return f, nil
	case reflect.String:
		tmp, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		l := int(binary.LittleEndian.Uint32(tmp))
		tmp2, err := read(r, l)
		if err != nil {
			return nil, err
		}
		s := string(tmp2)
		return s, nil
	case reflect.Array:
		l := t.Len()
		a := reflect.New(t).Elem()
		for i := 0; i < l; i++ {
			av, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			a.Index(i).Set(reflect.ValueOf(av))
		}
		return a.Interface(), nil
	case reflect.Slice:
		tmp, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		l := int(binary.LittleEndian.Uint32(tmp))
		a := reflect.New(t).Elem()
		for i := 0; i < l; i++ {
			av, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			a = reflect.Append(a, reflect.ValueOf(av))
		}
		return a.Interface(), nil
	case reflect.Map:
		tmp, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		l := int(binary.LittleEndian.Uint32(tmp))
		m := reflect.MakeMap(t)
		for i := 0; i < l; i++ {
			k, err := deserialize(t.Key(), r)
			if err != nil {
				return nil, err
			}
			v, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
		return m.Interface(), nil
	case reflect.Ptr:
		tmp, err := read(r, 1)
		if err != nil {
			return nil, err
		}
		valid := uint8(tmp[0])
		if valid == 0 {
			return nil, nil
		} else {
			p := reflect.New(t.Elem())
			deref, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			p.Elem().Set(reflect.ValueOf(deref))
			return p.Interface(), nil
		}
	case reflect.Struct:
		s, err := deserializeStruct(t, r)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return nil, nil
}

func deserializeStruct(t reflect.Type, r *bytes.Reader) (interface{}, error) {
	v := reflect.New(t).Elem()

	fieldMap := make(map[string]int)
	fields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		tag := field.Tag
		if tag.Get("borsh_skip") == "true" {
			continue
		}
		fieldMap[name] = i
		fields = append(fields, name)
	}
	sort.Strings(fields)
	for _, field := range fields {
		fv, err := deserialize(t.Field(fieldMap[field]).Type, r)
		if err != nil {
			return nil, err
		}
		v.Field(fieldMap[field]).Set(reflect.ValueOf(fv))
	}

	return v.Interface(), nil
}

func Serialize(s interface{}) ([]byte, error) {
	result := new(bytes.Buffer)

	err := serialize(reflect.ValueOf(s), result)
	return result.Bytes(), err
}

func serializeStruct(v reflect.Value, b *bytes.Buffer) error {
	t := v.Type()

	fieldMap := make(map[string]int)
	fields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("borsh_skip") == "true" {
			continue
		}
		fieldMap[field.Name] = i
		fields = append(fields, field.Name)
	}
	sort.Strings(fields)
	for _, field := range fields {
		err := serialize(v.Field(fieldMap[field]), b)
		if err != nil {
			return err
		}
	}
	return nil
}

func serialize(v reflect.Value, b *bytes.Buffer) error {
	switch v.Kind() {
	case reflect.Int8:
		b.WriteByte(byte(v.Interface().(int8)))
	case reflect.Int16:
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, uint16(v.Interface().(int16)))
		b.Write(tmp)
	case reflect.Int32:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(v.Interface().(int32)))
		b.Write(tmp)
	case reflect.Int64:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, uint64(v.Interface().(int64)))
		b.Write(tmp)
	case reflect.Int:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, uint64(v.Interface().(int)))
		b.Write(tmp)
	case reflect.Uint8:
		b.WriteByte(byte(v.Uint()))
	case reflect.Uint16:
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, v.Interface().(uint16))
		b.Write(tmp)
	case reflect.Uint32:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, v.Interface().(uint32))
		b.Write(tmp)
	case reflect.Uint64, reflect.Uint:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, v.Uint())
		b.Write(tmp)
	case reflect.Float32:
		tmp := make([]byte, 4)
		f := v.Float()
		if f == math.NaN() {
			return errors.New("NaN float value")
		}
		binary.LittleEndian.PutUint32(tmp, math.Float32bits(float32(f)))
		b.Write(tmp)
	case reflect.Float64:
		tmp := make([]byte, 8)
		f := v.Float()
		if f == math.NaN() {
			return errors.New("NaN float value")
		}
		binary.LittleEndian.PutUint64(tmp, math.Float64bits(f))
		b.Write(tmp)
	case reflect.String:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(len(v.String())))
		b.Write(tmp)
		b.Write([]byte(v.String()))
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			err := serialize(v.Index(i), b)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(v.Len()))
		b.Write(tmp)
		for i := 0; i < v.Len(); i++ {
			err := serialize(v.Index(i), b)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(v.Len()))
		b.Write(tmp)
		keys := v.MapKeys()
		sort.Slice(keys, vComp(keys))
		for _, k := range keys {
			serialize(k, b)
			serialize(v.MapIndex(k), b)
		}
	case reflect.Ptr:
		if v.IsNil() {
			b.WriteByte(0)
		} else {
			b.WriteByte(1)
			err := serialize(v.Elem(), b)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		err := serializeStruct(v, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func vComp(keys []reflect.Value) func(int, int) bool {
	return func(i int, j int) bool {
		a, b := keys[i], keys[j]
		if a.Kind() == reflect.Interface {
			a = a.Elem()
			b = b.Elem()
		}
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() < b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return a.Uint() < b.Uint()
		case reflect.String:
			return a.String() < b.String()
		}
		panic("unsupported key compare")
	}
}
