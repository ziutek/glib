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

func TestSignal(t *testing.T) {
	o := NewObject(TYPE_OBJECT, nil)

	s := NewSignal("sig1", TYPE_NONE, TYPE_POINTER, TYPE_GO_INT)
	t.Logf("Signal: %s", s)

	o.Connect(s, (*A).handler)

	a := A("test_signal")
	o.Emit(s, &a, 123)
}

func TestSignalName(t *testing.T) {
	o := NewObject(TYPE_OBJECT, nil)

	NewSignal("sig2", TYPE_NONE, TYPE_POINTER, TYPE_GO_INT)

	o.ConnectByName("sig2", (*A).handler)

	a := A("test_signal_name")
	o.EmitByName("sig2", &a, 456)
}

type A string

func (a *A) handler(i int) {
	fmt.Printf("handler: %d", i)
}

