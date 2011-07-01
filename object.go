package glib

/*
#include <stdlib.h>
#include "closure.h"

static inline
gulong go_signal_connect(GObject* inst, guint sig, GoClosure* cl) {
	return g_signal_connect_closure_by_id(
		(gpointer) inst,
		sig,
		0,
		(GClosure*) cl,
		TRUE
	);
}

static inline
void go_signal_emit(const GValue *inst_and_params, guint sig, GValue *ret) {
	return g_signal_emitv(
		inst_and_params,
		sig,
		0,
		ret
	);
}
*/
import "C"

import (
	"runtime"
)

type SigHandlerId C.gulong

type Object struct {
	obj *C.GObject
	cls map[SigHandlerId]*Closure
}

func (o *Object) Ref() *Object {
	if o.obj == nil {
		panic("Ref on nil object")
	}
	return &Object{(*C.GObject)(C.g_object_ref(C.gpointer(o.obj))), o.cls}
}

func (o *Object) Unref() {
	if o.obj == nil {
		panic("Unref on nil object")
	}
	C.g_object_unref(C.gpointer(o.obj))
	o.obj = nil
	o.cls = nil
}

func (o *Object) destroy() {
	if o.obj != nil {
		C.g_object_unref(C.gpointer(o.obj))
	}
}

func MapGObject(obj *C.GObject) (o *Object) {
	o = &Object{obj, make(map[SigHandlerId]*Closure)}
	runtime.SetFinalizer(o, (*Object).destroy)
	return
}

func MapGPointer(p C.gpointer) (o *Object) {
	return MapGObject((*C.GObject)(p))
}

func NewObject(t Type) *Object {
	return MapGPointer(C.g_object_newv(C.GType(t), 0, nil))
}

func (o *Object) Connect(sig Signal, cb_func interface{}) {
	c := NewClosure(cb_func)
	h := SigHandlerId(C.go_signal_connect(o.obj, C.guint(sig), c.cl))
	o.cls[h] = c
}

func (o *Object) Emit(sig Signal, args ...interface{}) interface{} {
	prms := make([]Value, len(args) + 1)
	prms[0] = *ValueOf(o.obj)
	for i, a := range args {
		prms[i+1] = *ValueOf(a)
	}
	ret := NewValue()
	C.go_signal_emit((*C.GValue)(&prms[0]), C.guint(sig), (*C.GValue)(ret))
	return ret.Get()
}
