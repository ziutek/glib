package glib

/*
#include <glib-object.h>
*/
import "C"

type Signal C.guint

func NewSignal(name string, ot Type, it ...Type) Signal {
	sig_name := NewString(name)
	defer sig_name.Free()
	var cit *C.GType
	if len(it) > 0 {
		cit = (*C.GType)(&it[0])
	}
	return Signal(C.g_signal_newv(
		sig_name.G(),
		TYPE_OBJECT.GType(),
		C.G_SIGNAL_RUN_LAST,
		nil,
		nil,
		nil,
		nil,
		ot.GType(),
		C.guint(len(it)),
		cit,
	))
}

func (s Signal) String() string {
	return C.GoString((*C.char)(C.g_signal_name(C.guint(s))))
}

func SignalLookup(n string, t Type) Signal {
	s := NewString(n)
	defer s.Free()
	return Signal(C.g_signal_lookup(s.G(), t.GType()))
}
