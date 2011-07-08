package glib

/*
#include <glib-object.h>
*/
import "C"

import (
	"reflect"
	"fmt"
)

type ValueGetter interface {
	Value() *Value
}

type Value C.GValue

func (v *Value) g() *C.GValue {
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
		C.g_value_reset(v.g())

	case reflect.Bool:
		C.g_value_set_boolean(v.g(), gBoolean(r.Bool()))

	case reflect.Int:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_int(v.g(), C.gint(i.(int)))
		} else {
			C.g_value_set_long(v.g(), C.glong(i.(int)))
		}

	case reflect.Int8:
		C.g_value_set_char(v.g(), C.gchar(i.(int8)))

	case reflect.Int32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_int(v.g(), C.gint(i.(int32)))
		} else {
			C.g_value_set_long(v.g(), C.glong(i.(int32)))
		}

	case reflect.Int64:
		C.g_value_set_int64(v.g(), C.gint64(i.(int64)))

	case reflect.Uint:
		if TYPE_GO_INT == TYPE_INT {
			C.g_value_set_uint(v.g(), C.guint(i.(uint)))
		} else {
			C.g_value_set_ulong(v.g(), C.gulong(i.(uint)))
		}

	case reflect.Uint8:
		C.g_value_set_uchar(v.g(), C.guchar(i.(uint8)))

	case reflect.Uint32:
		if TYPE_GO_INT32 == TYPE_INT {
			C.g_value_set_uint(v.g(), C.guint(i.(uint32)))
		} else {
			C.g_value_set_ulong(v.g(), C.gulong(i.(uint32)))
		}

	case reflect.Uint64:
		C.g_value_set_uint64(v.g(), C.guint64(i.(uint64)))

	case reflect.Float32:
		C.g_value_set_float(v.g(), C.gfloat(i.(float32)))

	case reflect.Float64:
		C.g_value_set_double(v.g(), C.gdouble(i.(float64)))

	case reflect.Ptr:
		C.g_value_set_pointer(v.g(), C.gpointer(r.Pointer()))

	case reflect.String:
		C.g_value_set_static_string(v.g(), (*C.gchar)(C.CString(r.String())))

	default:
		panic("&Value.Set Can't represent Go value in Glib type system.")
	}
}

// Initializes value with the default value of type. 
func (v *Value) Init(t Type) {
	C.g_value_init(v.g(), t.g())
}

// Clears the current value in value and "unsets" the type,
func (v *Value) Unset() {
	C.g_value_unset(v.g())
}

// Returns new initializes value
func DefaultValue(t Type) (v *Value) {
	v = new(Value)
	v.Init(t)
	return
}

// Returns a pointer to new Value initialized to the value stored in the
// interface i. If i contains pointer to Value returns this pointer.
func ValueOf(i interface{}) *Value {
	if v, ok := i.(*Value); ok {
		return v
	}
	v := DefaultValue(TypeOf(i))
	v.Set(i)
	return v
}

// Copies the value into dst.
func (v *Value) Copy(dst *Value) {
	C.g_value_copy(v.g(), dst.g())
}

func (v *Value) Get() interface{} {
	t := Type(v.g().g_type)
	switch t {
	case TYPE_INVALID:
		return nil

	case TYPE_STRING:
		return C.GoString((*C.char)(C.g_value_get_string(v.g())))

	case TYPE_GO_INT:
		if TYPE_GO_INT == TYPE_INT {
			return int(C.g_value_get_int(v.g()))
		} else {
			return int(C.g_value_get_long(v.g()))
		}

	case TYPE_GO_UINT:
		if TYPE_GO_INT == TYPE_INT {
			return uint(C.g_value_get_uint(v.g()))
		} else {
			return uint(C.g_value_get_ulong(v.g()))
		}

	case TYPE_CHAR:
		return int8(C.g_value_get_char(v.g()))

	case TYPE_UCHAR:
		return uint8(C.g_value_get_uchar(v.g()))

	case TYPE_GO_INT32:
		if TYPE_GO_INT32 == TYPE_INT {
			return int32(C.g_value_get_int(v.g()))
		} else {
			return int32(C.g_value_get_long(v.g()))
		}

	case TYPE_GO_UINT32:
		if TYPE_GO_INT32 == TYPE_INT {
			return uint32(C.g_value_get_uint(v.g()))
		} else {
			return uint32(C.g_value_get_ulong(v.g()))
		}

	case TYPE_INT64:
		return int64(C.g_value_get_int64(v.g()))

	case TYPE_UINT64:
		return uint64(C.g_value_get_uint64(v.g()))

	case TYPE_BOOLEAN:
		return (C.g_value_get_boolean(v.g()) != C.gboolean(0))

	case TYPE_FLOAT:
		return float32(C.g_value_get_float(v.g()))

	case TYPE_DOUBLE:
		return float64(C.g_value_get_double(v.g()))

	case TYPE_POINTER:
		return Pointer(C.g_value_get_pointer(v.g()))

	case TYPE_GTYPE:
		return Type(C.g_value_get_gtype(v.g()))
	}
	/*if t.IsA(TYPE_OBJECT) {
		o := new(Object)
		o.Set(Pointer(C.g_value_get_object(v.g())))
		return o
	}*/
	// Value of unknown type is returned as Pointer
	return Pointer(C.g_value_peek_pointer(v.g()))
}

func (v *Value) GetString() string {
	return v.Get().(string)
}

func (v *Value) GetInt() int {
	return v.Get().(int)
}

func (v *Value) GetUint() uint {
	return v.Get().(uint)
}

func (v *Value) GetInt8() int8 {
	return v.Get().(int8)
}

func (v *Value) GetUint8() uint8 {
	return v.Get().(uint8)
}

func (v *Value) GetInt32() int32 {
	return v.Get().(int32)
}

func (v *Value) GetUint32() uint32 {
	return v.Get().(uint32)
}

func (v *Value) GetInt64() int64 {
	return v.Get().(int64)
}

func (v *Value) GetUint64() uint64 {
	return v.Get().(uint64)
}

func (v *Value) GetBool() bool {
	return v.Get().(bool)
}

func (v *Value) GetFloat32() float32 {
	return v.Get().(float32)
}

func (v *Value) GetFloat64() float64 {
	return v.Get().(float64)
}

func (v *Value) GetPointer() Pointer {
	return v.Get().(Pointer)
}

func (v *Value) GetObject() *Object {
	return v.Get().(*Object)
}

func (v *Value) GetType() Type {
	return v.Get().(Type)
}

func (v *Value) String() string {
	return fmt.Sprint(v.Get())
}
