package borsh_test

import (
	"github.com/ouromoros/borsh-go"
	"math/rand"
	"reflect"
	"testing"
)

func TestFuzz(t *testing.T) {
	s1 := rand.NewSource(42)
	r1 := rand.New(s1)

	for i := 0; i < 1000; i++ {
		st := fuzzType(r1, 0)
		for j := 0; j < 10; j++ {
			val := fuzzValue(r1, st)
			testValue(t, val)
		}
	}
}

func testValue(t *testing.T, v interface{}) {
	data, err := borsh.Serialize(v)
	if err != nil {
		t.Error(err)
	}
	parsed := reflect.New(reflect.TypeOf(v))
	err = borsh.Deserialize(parsed.Interface(), data)
	if err != nil {
		t.Error(err)
	}
	reflect.DeepEqual(v, parsed.Elem().Interface())
}

func fuzzValue(r *rand.Rand, t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Int8:
		return int8(r.Int())
	case reflect.Int16:
		return int16(r.Int())
	case reflect.Int32:
		return r.Int31()
	case reflect.Int64:
		return r.Int63()
	case reflect.Uint8:
		return uint8(r.Int())
	case reflect.Uint16:
		return uint16(r.Uint32())
	case reflect.Uint32:
		return r.Uint32()
	case reflect.Uint64:
		return r.Uint64()
	case reflect.Float32:
		return r.Float32()
	case reflect.Float64:
		return r.Float64()
	case reflect.String:
		return randomString(r)
	case reflect.Array:
		l := t.Len()
		a := reflect.New(t).Elem()
		for i := 0; i < l; i++ {
			av := fuzzValue(r, t.Elem())
			a.Index(i).Set(reflect.ValueOf(av))
		}
		return a.Interface()
	case reflect.Slice:
		a := reflect.New(t).Elem()
		l := r.Int() % 10
		for i := 0; i < l; i++ {
			av := fuzzValue(r, t.Elem())
			a = reflect.Append(a, reflect.ValueOf(av))
		}
		return a.Interface()
	case reflect.Map:
		l := r.Int() % 10
		m := reflect.MakeMap(t)
		for i := 0; i < l; i++ {
			k := fuzzValue(r, t.Key())
			v := fuzzValue(r, t.Elem())
			m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
		return m.Interface()
	case reflect.Ptr:
		valid := r.Int()%2 == 1
		if valid {
			p := reflect.New(t.Elem())
			de := fuzzValue(r, t.Elem())
			p.Elem().Set(reflect.ValueOf(de))
			return p.Interface()
		} else {
			p := reflect.New(t.Elem())
			return p.Interface()
		}
	case reflect.Struct:
		v := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fv := fuzzValue(r, field.Type)
			v.Field(i).Set(reflect.ValueOf(fv))
		}
		return v.Interface()
	}
	panic(t)
}

func fuzzType(r *rand.Rand, c int) reflect.Type {
	p := new(uint)
	if c >= 10 {
		// prevent too deep nested structure
		return fuzzBasicType(r)
	}
	fuzzTypes := []reflect.Type{
		reflect.TypeOf([]int{}),
		reflect.TypeOf(map[int]int{}),
		reflect.TypeOf([1]int{}),
		reflect.TypeOf(p),
		reflect.TypeOf(struct{}{}),

		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
	}
	i := rand.Int() % len(fuzzTypes)
	switch fuzzTypes[i].Kind() {
	case reflect.Struct:
		nMembers := rand.Int() % 10
		fields := make([]reflect.StructField, 0)
		for i := 0; i < nMembers; i++ {
			fields = append(fields, randomField(r, c+1))
		}
		return reflect.StructOf(fields)
	case reflect.Map:
		k := fuzzBasicType(r)
		v := fuzzType(r, c+1)
		return reflect.MapOf(k, v)
	case reflect.Ptr:
		t := fuzzType(r, c+1)
		return reflect.PtrTo(t)
	case reflect.Slice:
		t := fuzzType(r, c+1)
		return reflect.SliceOf(t)
	case reflect.Array:
		l := r.Int() % 10
		t := fuzzType(r, c+1)
		return reflect.ArrayOf(l, t)
	default:
		return fuzzTypes[i]
	}
}

func fuzzBasicType(r *rand.Rand) reflect.Type {
	fuzzTypes := []reflect.Type{
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
	}
	i := r.Int() % len(fuzzTypes)
	return fuzzTypes[i]
}

func randomField(r *rand.Rand, c int) reflect.StructField {
	name := randomString(r)
	t := fuzzType(r, c+1)
	return reflect.StructField{
		Type: t,
		Name: name,
	}
}

func randomString(r *rand.Rand) string {
	s := []rune("ABCDEFGHIJKLMNOPQRSTUVW")
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return string(s)
}
