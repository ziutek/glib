package glib

import (
	"testing"
	"fmt"
)

// GLib values
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

// Signals by ID
func TestSignal(t *testing.T) {
	// Object on which signal will be emitted
	o := NewObject(TYPE_OBJECT, nil)

	// New signal without return value, which can be emitted on TYPE_OBJECT
	// and accepts one parameter of type TYPE_GO_INT
	s := NewSignal("sig1", TYPE_NONE, o.Type(), TYPE_GO_INT)
	t.Logf("Signal: %s", s)

	// Object that have methods which will be used as signal handlers
	a := A("test_signal")

	// Connect to a.handler method. o will be passed as first argument.
	o.ConnectSid(s, (*A).handler, &a)
	// Connect to a.noi_hr method. o wiln not be passed to method.
	o.ConnectSidNoi(s, (*A).noi_h, &a)
	// Connect to fh function. o will be passed as first argument.
	o.ConnectSid(s, fh, nil)

	// Emit signal with 123 as its TYPE_GO_INT argument
	o.EmitById(s, 123)
}

// Like TestSignal but uses signal names
func TestSignalName(t *testing.T) {
	o := NewObject(TYPE_OBJECT, nil)

	NewSignal("sig2", TYPE_NONE, TYPE_OBJECT, TYPE_GO_INT)

	a := A("test_signal_name")

	o.Connect("sig2", (*A).handler, &a)
	o.ConnectNoi("sig2", (*A).noi_h, &a)
	o.Connect("sig2", fh, nil)

	o.Emit("sig2", 456)
}

type A string

func (a *A) handler(o *Object, i int) {
	fmt.Printf("handler: %s, %v, %d\n", a, o, i)
}

func (a *A) noi_h(i int) {
	fmt.Printf("noi_h: %s, %d\n", a, i)
}

func fh(o *Object, i int) {
	fmt.Printf("fh: %v, %d\n", o, i)
}

