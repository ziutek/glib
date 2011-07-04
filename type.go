package glib

/*
#include <glib-object.h>

#define _GINT_SIZE sizeof(gint)
#define _GLONG_SIZE sizeof(glong)
*/
import "C"

import (
	"strconv"
	"reflect"
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

var (
	TYPE_GTYPE     Type
	TYPE_GO_INT    Type
	TYPE_GO_UINT   Type
	TYPE_GO_INT32  Type
	TYPE_GO_UINT32 Type
)

// Returns the Type of the value in the interface{}.
func TypeOf(i interface{}) Type {
	// Types defined in our package
	if o, ok := i.(ObjectI); ok {
		return o.Type()
	}
	if _, ok := i.(Type); ok {
		return TYPE_GTYPE
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

func (t Type) GType() C.GType {
	return C.GType(t)
}

func (t Type) String() string {
	return C.GoString((*C.char)(C.g_type_name(t.GType())))
}

var oi = reflect.TypeOf((*ObjectI)(nil)).Elem()

func (t Type) Match(rt reflect.Type) bool {
	if t == TYPE_OBJECT {
		return rt.Implements(oi)
	}
	k := rt.Kind()
	switch t {
	case TYPE_INVALID:
		return k == reflect.Invalid

	case TYPE_STRING:
		return k == reflect.String

	case TYPE_GO_INT:
		return k == reflect.Int

	case TYPE_GO_UINT:
		return k == reflect.Uint

	case TYPE_CHAR:
		return k == reflect.Int8

	case TYPE_UCHAR:
		return k == reflect.Uint8

	case TYPE_GO_INT32:
		return k == reflect.Int32

	case TYPE_GO_UINT32:
		return k == reflect.Uint32

	case TYPE_INT64:
		return k == reflect.Int64

	case TYPE_UINT64:
		return k == reflect.Uint64

	case TYPE_BOOLEAN:
		return k == reflect.Bool

	case TYPE_FLOAT:
		return k == reflect.Float32

	case TYPE_DOUBLE:
		return k == reflect.Float64

	case TYPE_POINTER:
		return k == reflect.Ptr
	}
	return false
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
