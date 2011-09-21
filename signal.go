package glib

/*
#include <glib-object.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
	"strings"
)

type SignalFlags C.GSignalFlags

const (
  SIGNAL_RUN_FIRST = C.G_SIGNAL_RUN_FIRST
  SIGNAL_RUN_LAST = C.G_SIGNAL_RUN_LAST
  SIGNAL_RUN_CLEANUP = C.G_SIGNAL_RUN_CLEANUP
  SIGNAL_NO_RECURSE = C.G_SIGNAL_NO_RECURSE
  SIGNAL_DETAILED = C.G_SIGNAL_DETAILED
  SIGNAL_ACTION = C.G_SIGNAL_ACTION
  SIGNAL_NO_HOOKS = C.G_SIGNAL_NO_HOOKS
)

type SignalId C.guint

func NewSignal(name string, rt, it Type, pt ...Type) SignalId {
	sig_name := C.CString(name)
	defer C.free(unsafe.Pointer(sig_name))
	var cit *C.GType
	if len(pt) > 0 {
		cit = (*C.GType)(&pt[0])
	}
	return SignalId(C.g_signal_newv(
		(*C.gchar)(sig_name),
		it.g(),
		SIGNAL_RUN_LAST,
		nil,
		nil,
		nil,
		nil,
		rt.g(),
		C.guint(len(pt)),
		cit,
	))
}

func (s SignalId) String() string {
	return C.GoString((*C.char)(C.g_signal_name(C.guint(s))))
}

func SignalLookup(name string, t Type) (sid SignalId, detail Quark) {
	st := strings.SplitN(name, "::", 2)
	if len(st) == 2 {
		detail = QuarkFromString(st[1])
	}
	s := C.CString(st[0])
	sid = SignalId(C.g_signal_lookup((*C.gchar)(s), t.g()))
	defer C.free(unsafe.Pointer(s))
	return
}
