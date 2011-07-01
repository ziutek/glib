package glib

/*
#include <stdlib.h>
#include <glib-object.h>

static inline guint go_signal_new(char* name, GType ot, guint ni, GType* it) {
	return g_signal_newv(
		name,
		G_TYPE_OBJECT,
		G_SIGNAL_RUN_LAST,
		NULL,
		NULL,
		NULL,
		NULL,
		ot,
		ni,
		it
	);
}
*/
import "C"

import (
	"unsafe"
)

type Signal C.guint

func NewSignal(name string, ot Type, it ...Type) Signal {
	sig_name := C.CString(name)
	defer C.free(unsafe.Pointer(sig_name))
	var cit *C.GType
	if len(it) > 0 {
		cit = (*C.GType)(&it[0])
	}
	return Signal(C.go_signal_new(sig_name, C.GType(ot), C.guint(len(it)), cit))
}

func (s Signal) String() string {
	return C.GoString((*C.char)(C.g_signal_name(C.guint(s))))
}
