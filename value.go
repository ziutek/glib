package glib

/*
#include <glib-object.h>
*/
import "C"

import (
	"reflect"
	"unsafe"
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
	r := reflect.ValueOf(i)
	switch r.Kind() {
	case reflect.Invalid:
		C.g_value_reset((*C.GValue)(v))

	case reflect.Bool:
		if r.Bool() {
			C.g_value_set_boolean((*C.GValue)(v), C.gboolean(1))
		} else {
			C.g_value_set_boolean((*C.GValue)(v), C.gboolean(0))
		}

	case reflect.Int:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_int((*C.GValue)(v), C.gint(i.(int)))
		} else {
			C.g_value_set_long((*C.GValue)(v), C.glong(i.(int)))
		}

	case reflect.Int8:
		C.g_value_set_char((*C.GValue)(v), C.gchar(i.(int8)))

	case reflect.Int32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_int((*C.GValue)(v), C.gint(i.(int32)))
		} else {
			C.g_value_set_long((*C.GValue)(v), C.glong(i.(int32)))
		}

	case reflect.Int64:
		C.g_value_set_int64((*C.GValue)(v), C.gint64(i.(int64)))

	case reflect.Uint:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_uint((*C.GValue)(v), C.guint(i.(uint)))
		} else {
			C.g_value_set_ulong((*C.GValue)(v), C.gulong(i.(uint)))
		}

	case reflect.Uint8:
		C.g_value_set_uchar((*C.GValue)(v), C.guchar(i.(uint8)))

	case reflect.Uint32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_uint((*C.GValue)(v), C.guint(i.(uint32)))
		} else {
			C.g_value_set_ulong((*C.GValue)(v), C.gulong(i.(uint32)))
		}

	case reflect.Uint64:
		C.g_value_set_uint64((*C.GValue)(v), C.guint64(i.(uint64)))

	case reflect.Float32:
		C.g_value_set_float((*C.GValue)(v), C.gfloat(i.(float32)))

	case reflect.Float64:
		C.g_value_set_double((*C.GValue)(v), C.gdouble(i.(float64)))

	case reflect.Ptr:
		if o, ok := i.(*Object); ok {
			C.g_value_set_object((*C.GValue)(v), C.gpointer(o))
		} else {
			C.g_value_set_pointer(
				(*C.GValue)(v),
				C.gpointer(unsafe.Pointer(r.Pointer())),
			)
		}

	case reflect.String:
		C.g_value_set_static_string(
			(*C.GValue)(v),
			(*C.gchar)(C.CString(r.String())),
		)

	case reflect.UnsafePointer:
		C.g_value_set_pointer((*C.GValue)(v), C.gpointer(i.(unsafe.Pointer)))

	default:
		panic("Can't represent Go value in Glib type system.")
	}
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
	v := new(Value)
	v.Init(t)
	return v
}

// Returns a pointer to new Value initialized to the value stored in the
// interface i. If i contains pointer to Value returns this pointer.
func ValueOf(i interface{}) *Value {
	if v, ok := i.(*Value); ok {
		return v
	}
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
	case TYPE_INVALID:
		return nil

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
		return (*Object)(C.g_value_get_object((*C.GValue)(v)))
	}
	panic("Can't represent GLib value in Go type system.")
}

func (v *Value) String() string {
	return fmt.Sprint(v.Get())
}
