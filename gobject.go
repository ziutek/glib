package glib

/*
#cgo pkg-config: glib-2.0 gobject-2.0

#include <glib-object.h>
*/
import "C"

import (
	"runtime"
)

// A numerical value which represents the unique identifier of a registered type
type Type C.GType

const (
	TYPE_INVALID   = Type(C.G_TYPE_INVALID)
	TYPE_NONE      = Type(C.G_TYPE_NONE)
	TYPE_INTERFACE = Type(C.G_TYPE_INTERFACE)
	TYPE_CHAR      = Type(C.G_TYPE_CHAR)
	TYPE_UCHAR     = Type(C.G_TYPE_UCHAR)
	TYPE_BOOLEAN   = Type(C.G_TYPE_BOOLEAN)
	TYPE_INT       = Type(C.G_TYPE_INT)
	TYPE_UINT      = Type(C.G_TYPE_UINT)
	TYPE_LONG      = Type(C.G_TYPE_LONG)
	TYPE_ULONG     = Type(C.G_TYPE_ULONG)
	TYPE_INT64     = Type(C.G_TYPE_INT64)
	TYPE_UINT64    = Type(C.G_TYPE_UINT64)
	TYPE_ENUM      = Type(C.G_TYPE_ENUM)
	TYPE_FLAGS     = Type(C.G_TYPE_FLAGS)
	TYPE_FLOAT     = Type(C.G_TYPE_FLOAT)
	TYPE_DOUBLE    = Type(C.G_TYPE_DOUBLE)
	TYPE_STRING    = Type(C.G_TYPE_STRING)
	TYPE_POINTER   = Type(C.G_TYPE_POINTER)
	TYPE_BOXED     = Type(C.G_TYPE_BOXED)
	TYPE_PARAM     = Type(C.G_TYPE_PARAM)
	TYPE_OBJECT    = Type(C.G_TYPE_OBJECT)
	TYPE_VARIANT   = Type(C.G_TYPE_VARIANT)
)

var TYPE_GTYPE = Type(C.g_gtype_get_type())


type Object struct {
	obj *C.GObject
}

func (o Object) Ref() Object {
	if o.obj == nil {
		panic("Ref on nil object")
	}
	return Object{(*C.GObject)(C.g_object_ref(C.gpointer(o.obj)))}
}

func (o Object) Unref() {
	if o.obj == nil {
		panic("Unref on nil object")
	}
	C.g_object_unref(C.gpointer(o.obj))
	o.obj = nil
}

func (o Object) destroy() {
	if o.obj != nil {
		C.g_object_unref(C.gpointer(o.obj))
	}
}

func MapGObject(obj *C.GObject) (o Object) {
	o.obj = obj
	runtime.SetFinalizer(o, (*Object).destroy)
	return
}

func MapGPointer(p C.gpointer) (o Object) {
	return MapGObject((*C.GObject)(p))
}

func NewObject(t Type) Object {
	return MapGPointer(C.g_object_newv(C.GType(t), 0, nil))
}

func (o Object) Emit(name string, ...interface{}) {

}

func init() {
	C.g_thread_init(nil)
	C.g_type_init()
}
