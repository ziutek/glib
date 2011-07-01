package glib

/*
#include <glib-object.h>
*/
import "C"

import (
	"unsafe"
	"runtime"
	"fmt"
)

// An opaque structure used to hold different types of values.
type Value C.GValue

// Returns v's type.
func (v *Value) Type() Type {
	return Type(v.g_type)
}

// Set value to i
func (v *Value) Set(i interface{}) {
	switch x := i.(type) {
	case string:
		C.g_value_set_static_string(
			(*C.GValue)(v),
			(*C.gchar)(C.CString(x)),
		)
	case int:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_int((*C.GValue)(v), C.gint(x))
		} else {
			C.g_value_set_long((*C.GValue)(v), C.glong(x))
		}
	case uint:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_uint((*C.GValue)(v), C.guint(x))
		} else {
			C.g_value_set_ulong((*C.GValue)(v), C.gulong(x))
		}
	case int8:
		C.g_value_set_char((*C.GValue)(v), C.gchar(x))
	case uint8:
		C.g_value_set_uchar((*C.GValue)(v), C.guchar(x))
	case int32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_int((*C.GValue)(v), C.gint(x))
		} else {
			C.g_value_set_long((*C.GValue)(v), C.glong(x))
		}
	case uint32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_uint((*C.GValue)(v), C.guint(x))
		} else {
			C.g_value_set_ulong((*C.GValue)(v), C.gulong(x))
		}
	case int64:
		C.g_value_set_int64((*C.GValue)(v), C.gint64(x))
	case uint64:
		C.g_value_set_uint64((*C.GValue)(v), C.guint64(x))
	case bool:
		if x {
			C.g_value_set_boolean((*C.GValue)(v), C.gboolean(1))
		} else {
			C.g_value_set_boolean((*C.GValue)(v), C.gboolean(0))
		}
	case float32:
		C.g_value_set_float((*C.GValue)(v), C.gfloat(x))
	case float64:
		C.g_value_set_double((*C.GValue)(v), C.gdouble(x))
	case unsafe.Pointer:
		C.g_value_set_pointer((*C.GValue)(v), C.gpointer(x))
	case *Object:
		C.g_value_set_object((*C.GValue)(v), C.gpointer(x.obj))
	case nil:
		C.g_value_reset((*C.GValue)(v))
	default:
		panic("Unknown type")
	}
}

// Returns new uninitialized value
func NewValue() *Value {
	v := new(Value)
	runtime.SetFinalizer(v, (*Value).Unset)
	return v
}

// Initializes value with the default value of type. 
func (v *Value) Init(t Type) {
	C.g_value_init((*C.GValue)(v), C.GType(t))
}

// Clears the current value in value and "unsets" the type,
func (v *Value) Unset() {
	C.g_value_unset((*C.GValue)(v))
}

// Returns new initializes value
func NewValueInit(t Type) *Value {
	v := NewValue()
	v.Init(t)
	return v
}

// Returns a new Value initialized to the value stored in the interface i.
func ValueOf(i interface{}) *Value {
	v := NewValueInit(TypeOf(i))
	v.Set(i)
	return v
}

// Copies the value into dst.
func (v *Value) Copy(dst *Value) {
	C.g_value_copy((*C.GValue)(v), (*C.GValue)(dst))
}

func (v *Value) Get() interface{} {
	switch Type((*C.GValue)(v).g_type) {
	case TYPE_STRING:
		return C.GoString((*C.char)(C.g_value_get_string((*C.GValue)(v))))
	case TYPE_GO_INT:
		if TYPE_GO_INT == TYPE_INT {
			return int(C.g_value_get_int((*C.GValue)(v)))
		} else {
			return int(C.g_value_get_long((*C.GValue)(v)))
		}
	case TYPE_GO_UINT:
		if TYPE_GO_INT == TYPE_INT {
			return uint(C.g_value_get_uint((*C.GValue)(v)))
		} else {
			return uint(C.g_value_get_ulong((*C.GValue)(v)))
		}
	case TYPE_CHAR:
		return int8(C.g_value_get_char((*C.GValue)(v)))
	case TYPE_UCHAR:
		return uint8(C.g_value_get_uchar((*C.GValue)(v)))
	case TYPE_GO_INT32:
		if TYPE_GO_INT32 == TYPE_INT {
			return int32(C.g_value_get_int((*C.GValue)(v)))
		} else {
			return int32(C.g_value_get_long((*C.GValue)(v)))
		}
	case TYPE_GO_UINT32:
		if TYPE_GO_INT32 == TYPE_INT {
			return uint32(C.g_value_get_uint((*C.GValue)(v)))
		} else {
			return uint32(C.g_value_get_ulong((*C.GValue)(v)))
		}
	case TYPE_INT64:
		return int64(C.g_value_get_int64((*C.GValue)(v)))
	case TYPE_UINT64:
		return uint64(C.g_value_get_uint64((*C.GValue)(v)))
	case TYPE_BOOLEAN:
		return (C.g_value_get_boolean((*C.GValue)(v)) != C.gboolean(0))
	case TYPE_FLOAT:
		return float32(C.g_value_get_float((*C.GValue)(v)))
	case TYPE_DOUBLE:
		return float64(C.g_value_get_double((*C.GValue)(v)))
	case TYPE_POINTER:
		return unsafe.Pointer(C.g_value_get_pointer((*C.GValue)(v)))
	case TYPE_OBJECT:
		return MapGPointer(C.g_value_get_object((*C.GValue)(v)))
	}
	// TODO - shoulda work with more GLIB types
	panic("Can't represent value in Go type system.")
}

func (v *Value) String() string {
	return fmt.Sprint(v.Get())
}
