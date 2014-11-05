// Bindings for glib
package glib

/*
#include <stdlib.h>
#include <glib-object.h>

#cgo CFLAGS: -Wno-deprecated-declarations
#cgo pkg-config: glib-2.0 gobject-2.0 gthread-2.0
*/
import "C"

import (
	"reflect"
	"unsafe"
)

type TypeGetter interface {
	Type() Type
}

type PointerSetter interface {
	SetPtr(p Pointer)
}

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
	TYPE_GO_INT32  = TYPE_INT
	TYPE_UINT      = Type(C.G_TYPE_UINT)
	TYPE_GO_UINT32 = TYPE_UINT
	TYPE_LONG      = Type(C.G_TYPE_LONG)
	TYPE_GO_INT    = TYPE_LONG
	TYPE_ULONG     = Type(C.G_TYPE_ULONG)
	TYPE_GO_UINT   = TYPE_ULONG
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

var TYPE_GTYPE Type

func (t Type) g() C.GType {
	return C.GType(t)
}

func (t Type) String() string {
	return C.GoString((*C.char)(C.g_type_name(t.g())))
}

func (t Type) QName() Quark {
	return Quark(C.g_type_qname(t.g()))
}

func (t Type) Type() Type {
	return TYPE_GTYPE
}

func (t Type) Value() *Value {
	v := NewValue(t.Type())
	C.g_value_set_gtype(v.g(), t.g())
	return v
}

func (t Type) Parent() Type {
	return Type(C.g_type_parent(t.g()))
}

func (t Type) Depth() int {
	return int(C.g_type_depth(t.g()))
}

// Returns the type that is derived directly from root type which is also
// a base class of t
func (t Type) NextBase(root Type) Type {
	return Type(C.g_type_next_base(t.g(), root.g()))
}

// If t type is a derivable type, check whether type is a descendant of
// it type. If t type is an glib interface, check whether type conforms
// to it.
func (t Type) IsA(it Type) bool {
	return C.g_type_is_a(t.g(), it.g()) != 0
}

var type_getter = reflect.TypeOf((*TypeGetter)(nil)).Elem()
var object_caster = reflect.TypeOf((*ObjectCaster)(nil)).Elem()

func (t Type) Match(rt reflect.Type) bool {
	if rt.Implements(object_caster) {
		return t.IsA(TYPE_OBJECT)
	}
	if rt.Implements(type_getter) {
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		r := reflect.New(rt).Interface().(TypeGetter).Type()
		return t.QName() == r.QName()
	}
	switch rt.Kind() {
	case reflect.Invalid:
		return t == TYPE_INVALID

	case reflect.String:
		return t == TYPE_STRING

	case reflect.Int:
		return t == TYPE_GO_INT

	case reflect.Uint:
		return t == TYPE_GO_UINT

	case reflect.Int8:
		return t == TYPE_CHAR

	case reflect.Uint8:
		return t == TYPE_UCHAR

	case reflect.Int32:
		return t == TYPE_GO_INT32

	case reflect.Uint32:
		return t == TYPE_GO_UINT32

	case reflect.Int64:
		return t == TYPE_INT64

	case reflect.Uint64:
		return t == TYPE_UINT64

	case reflect.Bool:
		return t == TYPE_BOOLEAN

	case reflect.Float32:
		return t == TYPE_FLOAT

	case reflect.Float64:
		return t == TYPE_DOUBLE

	case reflect.Ptr:
		return t == TYPE_POINTER
	}
	return false
}

func (t Type) Compatible(dst Type) bool {
	return C.g_value_type_compatible(t.g(), dst.g()) != C.gboolean(0)
}

func (t Type) Transformable(dst Type) bool {
	return C.g_value_type_transformable(t.g(), dst.g()) != C.gboolean(0)
}

// TypeOf returns the Type of the value in the i}.
func TypeOf(i interface{}) Type {
	// Types ov values that implements TypeGetter
	if o, ok := i.(TypeGetter); ok {
		return o.Type()
	}
	// Other types
	switch reflect.TypeOf(i).Kind() {
	case reflect.Invalid:
		return TYPE_INVALID

	case reflect.Bool:
		return TYPE_BOOLEAN

	case reflect.Int:
		return TYPE_GO_INT

	case reflect.Int8:
		return TYPE_CHAR

	case reflect.Int32:
		return TYPE_GO_INT32

	case reflect.Int64:
		return TYPE_INT64

	case reflect.Uint:
		return TYPE_GO_UINT

	case reflect.Uint8:
		return TYPE_UCHAR

	case reflect.Uint32:
		return TYPE_GO_UINT32

	case reflect.Uint64:
		return TYPE_UINT64

	case reflect.Float32:
		return TYPE_FLOAT

	case reflect.Float64:
		return TYPE_DOUBLE

	case reflect.Ptr:
		return TYPE_POINTER

	case reflect.String:
		return TYPE_STRING
	}
	panic("Can't map Go type to Glib type")
}

func TypeFromName(name string) Type {
	tn := C.CString(name)
	defer C.free(unsafe.Pointer(tn))
	return Type(C.g_type_from_name((*C.gchar)(tn)))
}

func init() {
	C.g_type_init()
	TYPE_GTYPE = Type(C.g_gtype_get_type())
}

type Pointer C.gpointer

func gBoolean(b bool) C.gboolean {
	if b {
		return C.TRUE
	}
	return C.FALSE
}

type Quark C.GQuark

func (q Quark) GQuark() C.GQuark {
	return C.GQuark(q)
}

func (q Quark) String() string {
	return C.GoString((*C.char)(C.g_quark_to_string(q.GQuark())))
}

func QuarkFromString(s string) Quark {
	return Quark(C.g_quark_from_static_string((*C.gchar)(C.CString(s))))
}

type Error C.GError

func (e *Error) Error() string {
	return C.GoString((*C.char)(e.message))
}

func (e *Error) GetDomain() Quark {
	return Quark(e.domain)
}

func (e *Error) GetCode() int {
	return int(e.code)
}
