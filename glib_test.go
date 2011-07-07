package glib

import (
	"testing"
	"fmt"
)

func TestValue(t *testing.T) {
	v1 := uint64(0xdeadbeaf)
	a := ValueOf(v1)
	b := DefaultValue(TYPE_UINT64)
	a.Copy(b)
	t.Logf("a = %s(%s), b = %s(%s)", a.Type(), a, b.Type(), b)
	if b.Get().(uint64) != v1 {
		t.Error("TYPE_UINT64")
	}
	v2 := -1
	a = ValueOf(v2)
	b = DefaultValue(TYPE_INT)
	a.Copy(b)
	t.Logf("a = %s(%s), b = %s(%s)", a.Type(), a, b.Type(), b)
	if b.Get() != v2 {
		t.Error("TYPE_INT")
	}
}

type A string

func (a *A) handler(o *Object, i int) {
	fmt.Printf("handler: %s, %v, %d\n", a, o, i)
}

func fh(o *Object, i int) {
	fmt.Printf("fh: %v, %d\n", o, i)
}

func TestSignal(t *testing.T) {
	o := NewObject(TYPE_OBJECT, nil)

	s := NewSignal("sig1", TYPE_NONE, TYPE_OBJECT, TYPE_GO_INT)
	t.Logf("Signal: %s", s)

	a := A("test_signal")

	o.ConnectById(s, (*A).handler, &a)
	o.ConnectById(s, fh, nil)

	o.EmitById(s, 123)
}

func TestSignalName(t *testing.T) {
	o := NewObject(TYPE_OBJECT, nil)

	NewSignal("sig2", TYPE_NONE, TYPE_OBJECT, TYPE_GO_INT)

	a := A("test_signal_name")

	o.Connect("sig2", (*A).handler, &a)
	o.Connect("sig2", fh, nil)

	o.Emit("sig2", 456)
}


