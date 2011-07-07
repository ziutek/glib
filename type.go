package glib

/*
#include <stdlib.h>
#include <glib-object.h>

#define _GINT_SIZE sizeof(gint)
#define _GLONG_SIZE sizeof(glong)

#cgo pkg-config: glib-2.0 gobject-2.0
*/
import "C"

import (
	"strconv"
	"reflect"
)

type TypeGetter interface {
	Type() Type
}

type PointerSetter interface {
	Set(p Pointer)
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

var (
	TYPE_GTYPE     Type
	TYPE_GO_INT    Type
	TYPE_GO_UINT   Type
	TYPE_GO_INT32  Type
	TYPE_GO_UINT32 Type
)


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
	v := DefaultValue(t.Type())
	C.g_value_set_gtype(v.g(), t.g())
	return v
}

func (t Type) Parent() Type {
	return Type(C.g_type_parent(t.g()))
}

func (t Type) Depth() uint {
	return uint(C.g_type_depth(t.g()))
}

// Returns the type that is derived directly from root type which is also
// a base class of t
func (t Type) NextBase(root Type) Type {
	return Type(C.g_type_next_base(t.g(), root.g()))
}

// If is_a type is a derivable type, check whether type is a descendant of
// is_a type. If is_a type is an interface, check whether type conforms to it.
func (t Type) IsA(it Type) bool {
	return C.g_type_is_a(t.g(), it.g()) != 0
}

var tg = reflect.TypeOf((*TypeGetter)(nil)).Elem()
var og = reflect.TypeOf((*ObjectGetter)(nil)).Elem()

func (t Type) Match(rt reflect.Type) bool {
	if rt.Implements(og) {
		return t.IsA(TYPE_OBJECT)
	}
	if rt.Implements(tg) {
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

// Returns the Type of the value in the interface{}.
func TypeOf(i interface{}) Type {
	// Types defined in our package
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

func init() {
	C.g_thread_init(nil)
	C.g_type_init()
	TYPE_GTYPE = Type(C.g_gtype_get_type())
	int_bytes := strconv.IntSize / 8
	if int_bytes == uint(C._GINT_SIZE) {
		TYPE_GO_INT = TYPE_INT
		TYPE_GO_UINT = TYPE_UINT
	} else if int_bytes == C._GLONG_SIZE {
		TYPE_GO_INT = TYPE_LONG
		TYPE_GO_UINT = TYPE_ULONG
	} else if int_bytes == 64 {
		TYPE_GO_INT = TYPE_INT64
		TYPE_GO_UINT = TYPE_UINT64
	} else {
		panic("Unexpectd size of 'int'")
	}
	int32_bytes := C.uint(4)
	if int32_bytes == C._GINT_SIZE {
		TYPE_GO_INT32 = TYPE_INT
		TYPE_GO_UINT32 = TYPE_UINT
	} else if int32_bytes == C._GLONG_SIZE {
		TYPE_GO_INT32 = TYPE_LONG
		TYPE_GO_UINT32 = TYPE_ULONG
	} else {
		panic("Neither gint nor glong are 32 bit numbers")
	}
}

/*type String []int8

func NewString(s string) String {
	return (*[1<<31-1]C.gchar)(unsafe.Pointer(C.CString(s)))[:len(s)+1]
}

func (s String) Ptr() Pointer {
	return Pointer(&s[0])
}

func (s String) Free() {
	C.free(unsafe.Pointer(s.Ptr()))
}*/

type Pointer C.gpointer

func GBoolean(b bool) (r C.gboolean) {
	if b {
		r = 1
	}
	return
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
