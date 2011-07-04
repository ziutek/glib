package glib

import (
	"testing"
	"fmt"
)

func TestValue(t *testing.T) {
	v1 := uint64(0xdeadbeaf)
	a := ValueOf(v1)
	b := NewValueInit(TYPE_UINT64)
	a.Copy(b)
	t.Logf("a = %s(%s), b = %s(%s)", a.Type(), a, b.Type(), b)
	if b.Get().(uint64) != v1 {
		t.Error("TYPE_UINT64")
	}
	v2 := -1
	a = ValueOf(v2)
	b = NewValueInit(TYPE_INT)
	a.Copy(b)
	t.Logf("a = %s(%s), b = %s(%s)", a.Type(), a, b.Type(), b)
	if b.Get() != v2 {
		t.Error("TYPE_INT")
	}
}

func TestObject(t *testing.T) {
	o := NewObj(TYPE_OBJECT, nil)

	s := NewSignal("bla", TYPE_NONE, TYPE_POINTER, TYPE_GO_INT)
	t.Logf("Signal: %s", s)

	Connect(o, s, (*A).handler)

	a := A("babababab")
	Emit(o, s, &a, 123)
}

type A string

func (a *A) handler(i int) {
	fmt.Printf("handler: %d", i)
}

