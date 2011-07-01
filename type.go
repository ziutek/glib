package glib

/*
#include <glib-object.h>

#define _GINT_SIZE sizeof(gint)
#define _GLONG_SIZE sizeof(glong)
*/
import "C"

import (
	"strconv"
	"unsafe"
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
	switch i.(type) {
	case string:
		return TYPE_STRING
	case int:
		return TYPE_GO_INT
	case uint:
		return TYPE_GO_UINT
	case int8:
		return TYPE_CHAR
	case uint8:
		return TYPE_UCHAR
	case int32:
		return TYPE_GO_INT32
	case uint32:
		return TYPE_GO_UINT32
	case int64:
		return TYPE_INT64
	case uint64:
		return TYPE_UINT64
	case bool:
		return TYPE_BOOLEAN
	case float32:
		return TYPE_FLOAT
	case float64:
		return TYPE_DOUBLE
	case unsafe.Pointer:
		return TYPE_POINTER
	case *Object:
		return TYPE_OBJECT
	case Type:
		return TYPE_GTYPE
	}
	return TYPE_INVALID
}

func (t Type) String() string {
	return C.GoString((*C.char)(C.g_type_name(C.GType(t))))
}

func init() {
	C.g_thread_init(nil)
	C.g_type_init()
	TYPE_GTYPE = Type(C.g_gtype_get_type())
	switch strconv.IntSize / 8 {
	case uint(C._GINT_SIZE):
		TYPE_GO_INT = TYPE_INT
		TYPE_GO_UINT = TYPE_UINT
	case uint(C._GLONG_SIZE):
		TYPE_GO_INT = TYPE_LONG
		TYPE_GO_UINT = TYPE_ULONG
	case 64:
		TYPE_GO_INT = TYPE_INT64
		TYPE_GO_UINT = TYPE_UINT64
	default:
		panic("Unexpectd size of 'int'")
	}
	switch C.uint(4) {
	case C._GINT_SIZE:
		TYPE_GO_INT32 = TYPE_INT
		TYPE_GO_UINT32 = TYPE_UINT
	case C._GLONG_SIZE:
		TYPE_GO_INT32 = TYPE_LONG
		TYPE_GO_UINT32 = TYPE_ULONG
	default:
		panic("Neither gint nor glong are 32 bit numbers")
	}
}
