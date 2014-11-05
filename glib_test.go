package glib

import (
	"fmt"
	"testing"
)

// GLib values
func TestValue(t *testing.T) {
	v1 := uint64(0xdeadbeaf)
	a := ValueOf(v1)
	b := NewValue(TYPE_UINT64)
	a.Copy(b)
	t.Logf("a = %s(%s), b = %s(%s)", a.Type(), a, b.Type(), b)
	if b.Get().(uint64) != v1 {
		t.Error("TYPE_UINT64")
	}
	v2 := -1
	a = ValueOf(v2)             // TYPE_GO_INT == TYPE_LONG
	b = NewValue(TYPE_GO_INT32) // TYPE_GO_INT32 == TYPE_INT
	a.Transform(b)
	t.Logf("a = %s(%s), b = %s(%s)", a.Type(), a, b.Type(), b)
	if b.Get().(int32) != int32(v2) {
		t.Error("TYPE_GO_INT32 (TYPE_INT)")
	}
}

// Signals by ID
func TestSignal(t *testing.T) {
	// An object on which the signal will be emitted
	o := NewObject(TYPE_OBJECT, nil)

	// The new signal without return value, which can be emitted on TYPE_OBJECT
	// and accepts one parameter of type TYPE_GO_INT
	s := NewSignal("sig1", TYPE_NONE, o.Type(), TYPE_GO_INT)
	t.Logf("Signal: %s", s)

	// The Go variable of type A that have some methods. These methods will be
	// used as signal handlers.
	a := A("Test with IDs")

	// Connect a.handler(*Object, int) to the signal.
	// o will be passed as its first argument. Second argument of type int will
	// be passed from second argument passed to the Emit function.
	o.ConnectSid(s, 0, (*A).handler, &a)

	// Connect a.noiHandler(int) to the signal.
	// o will not be passed to the method. An argument of type int will be
	// passed from second argument passed to the Emit function.
	o.ConnectSidNoi(s, 0, (*A).noiHandler, &a)

	// Connect funcHandler(*Object, int)to the signal.
	// o will be passed as its first argument. Second argument of type int
	// will be passed from second argument passed to the Emit function.
	o.ConnectSid(s, 0, funcHandler, nil)

	// Connect funcNoiHandler(int) to the signal.
	// o will not be passed to the function. An argument of type int will be
	// passed from second argument passed to the Emit function.
	o.ConnectSidNoi(s, 0, funcNoiHandler, nil)

	// Connect funcHandlerParam0(A, *Object, int) to the signal.
	// &a will be passed as its first argument, o will be passed as its second
	// argument. The thrid argument of type int will be from second argument
	// passed to the Emit function.
	o.ConnectSid(s, 0, funcHandlerParam0, &a)

	// Emit signal with 123 integer as argument.
	o.EmitById(s, 0, 123)
}

// Like TestSignal but uses signal names
func TestSignalName(t *testing.T) {
	o := NewObject(TYPE_OBJECT, nil)

	NewSignal("sig2", TYPE_NONE, TYPE_OBJECT, TYPE_GO_INT)

	a := A("Test with names")

	o.Connect("sig2", (*A).handler, &a)
	o.ConnectNoi("sig2", (*A).noiHandler, &a)
	o.Connect("sig2", funcHandler, nil)
	o.ConnectNoi("sig2", funcNoiHandler, nil)
	o.Connect("sig2", funcHandlerParam0, &a)

	o.Emit("sig2", 456)
}

type A string

func (a *A) handler(o *Object, i int) {
	fmt.Printf("handler: %s, %v, %d\n", *a, o, i)
}

func (a *A) noiHandler(i int) {
	fmt.Printf("noiHandler: %s, %d\n", *a, i)
}

func funcHandler(o *Object, i int) {
	fmt.Printf("funcHandler: %v, %d\n", o, i)
}

func funcNoiHandler(i int) {
	fmt.Printf("funcNoiHandler: %d\n", i)
}

func funcHandlerParam0(a *A, o *Object, i int) {
	fmt.Printf("funcHandlerParam0: %s %v %d\n", *a, o, i)
}
