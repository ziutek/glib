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

type ValueGetter interface {
	Value() *Value
}

// An opaque structure used to hold different types of values.
type Value C.GValue

func (v *Value) GValue() *C.GValue {
	return (*C.GValue)(v)
}

// Returns v's type.
func (v *Value) Type() Type {
	return Type(v.g_type)
}

// Set value to i
func (v *Value) Set(i interface{}) {
	if vg, ok := i.(ValueGetter); ok {
		vg.Value().Copy(v)
		return
	}
	// Other types
	r := reflect.ValueOf(i)
	switch r.Kind() {
	case reflect.Invalid:
		C.g_value_reset(v.GValue())

	case reflect.Bool:
		C.g_value_set_boolean(v.GValue(), GBoolean(r.Bool()))

	case reflect.Int:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_int(v.GValue(), C.gint(i.(int)))
		} else {
			C.g_value_set_long(v.GValue(), C.glong(i.(int)))
		}

	case reflect.Int8:
		C.g_value_set_char(v.GValue(), C.gchar(i.(int8)))

	case reflect.Int32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_int(v.GValue(), C.gint(i.(int32)))
		} else {
			C.g_value_set_long(v.GValue(), C.glong(i.(int32)))
		}

	case reflect.Int64:
		C.g_value_set_int64(v.GValue(), C.gint64(i.(int64)))

	case reflect.Uint:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_uint(v.GValue(), C.guint(i.(uint)))
		} else {
			C.g_value_set_ulong(v.GValue(), C.gulong(i.(uint)))
		}

	case reflect.Uint8:
		C.g_value_set_uchar(v.GValue(), C.guchar(i.(uint8)))

	case reflect.Uint32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_uint(v.GValue(), C.guint(i.(uint32)))
		} else {
			C.g_value_set_ulong(v.GValue(), C.gulong(i.(uint32)))
		}

	case reflect.Uint64:
		C.g_value_set_uint64(v.GValue(), C.guint64(i.(uint64)))

	case reflect.Float32:
		C.g_value_set_float(v.GValue(), C.gfloat(i.(float32)))

	case reflect.Float64:
		C.g_value_set_double(v.GValue(), C.gdouble(i.(float64)))

	case reflect.Ptr:
		C.g_value_set_pointer(v.GValue(), C.gpointer(r.Pointer()))

	case reflect.String:
		C.g_value_set_static_string(
			v.GValue(),
			NewString(r.String()).G(),
		)

	default:
		panic("Can't represent Go value in Glib type system.")
	}
}

// Initializes value with the default value of type. 
func (v *Value) Init(t Type) {
	C.g_value_init(v.GValue(), t.GType())
}

// Clears the current value in value and "unsets" the type,
func (v *Value) Unset() {
	C.g_value_unset(v.GValue())
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
	C.g_value_copy(v.GValue(), dst.GValue())
}

func (v *Value) Get() interface{} {
	switch Type(v.GValue().g_type) {
	case TYPE_INVALID:
		return nil

	case TYPE_STRING:
		return C.GoString((*C.char)(C.g_value_get_string(v.GValue())))

	case TYPE_GO_INT:
		if TYPE_GO_INT == TYPE_INT {
			return int(C.g_value_get_int(v.GValue()))
		} else {
			return int(C.g_value_get_long(v.GValue()))
		}

	case TYPE_GO_UINT:
		if TYPE_GO_INT == TYPE_INT {
			return uint(C.g_value_get_uint(v.GValue()))
		} else {
			return uint(C.g_value_get_ulong(v.GValue()))
		}

	case TYPE_CHAR:
		return int8(C.g_value_get_char(v.GValue()))

	case TYPE_UCHAR:
		return uint8(C.g_value_get_uchar(v.GValue()))

	case TYPE_GO_INT32:
		if TYPE_GO_INT32 == TYPE_INT {
			return int32(C.g_value_get_int(v.GValue()))
		} else {
			return int32(C.g_value_get_long(v.GValue()))
		}

	case TYPE_GO_UINT32:
		if TYPE_GO_INT32 == TYPE_INT {
			return uint32(C.g_value_get_uint(v.GValue()))
		} else {
			return uint32(C.g_value_get_ulong(v.GValue()))
		}

	case TYPE_INT64:
		return int64(C.g_value_get_int64(v.GValue()))

	case TYPE_UINT64:
		return uint64(C.g_value_get_uint64(v.GValue()))

	case TYPE_BOOLEAN:
		return (C.g_value_get_boolean(v.GValue()) != C.gboolean(0))

	case TYPE_FLOAT:
		return float32(C.g_value_get_float(v.GValue()))

	case TYPE_DOUBLE:
		return float64(C.g_value_get_double(v.GValue()))

	case TYPE_POINTER:
		return unsafe.Pointer(C.g_value_get_pointer(v.GValue()))

	case TYPE_OBJECT:
		o := new(Object)
		o.Set(Pointer(C.g_value_get_object(v.GValue())))
		return o

	case TYPE_GTYPE:
		return Type(C.g_value_get_gtype(v.GValue()))
	}
	panic("Unknown value type")
}

func (v *Value) String() string {
	return fmt.Sprint(v.Get())
}
