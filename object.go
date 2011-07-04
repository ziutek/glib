package glib

/*
#include <stdlib.h>
#include <glib-object.h>

typedef struct {
	GClosure cl;
	gpointer o; 
	gulong h_id;
} GoClosure;

typedef struct {
	GoClosure *cl;
	GValue *ret_val;
	guint n_param;
	const GValue *params;
	gpointer ih;
	gpointer data;
} MarshalParams;

extern void go_marshal(gpointer mp);  

static inline
void closure_marshal(GClosure* cl, GValue* ret_val, guint n_param,
		const GValue* params, gpointer ih, gpointer data) {
	MarshalParams mp = {(GoClosure*) cl, ret_val, n_param, params, ih, data};
	go_marshal(&mp);	
}

static inline
GoClosure* _closure_new(GObject *o) {
	GoClosure *cl = (GoClosure*) g_closure_new_simple(sizeof (GoClosure), NULL);
	cl->o = o;
	g_closure_set_marshal((GClosure *) cl, closure_marshal);
	return cl;
}

static inline
gulong _signal_connect(GObject* inst, guint sig, GoClosure* cl) {
	return g_signal_connect_closure_by_id(
		inst,
		sig,
		0,
		(GClosure*) cl,
		TRUE
	);
}

static inline
void _signal_emit(const GValue *inst_and_params, guint sig, GValue *ret) {
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
	"reflect"
	"unsafe"
	"fmt"
)

type SigHandlerId C.gulong


// Binding for any type derived from GObject should implement this interface.
type Object interface {
	// Should return self mapped to C.GObject pointer
	Obj() *C.GObject

	// Should return C.gpointer mapped to pointer to self type
	FromPtr(C.gpointer) Object
}

// Returns C pointer
func Ref(o Object) Object {
	return o.FromPtr(C.g_object_ref(C.gpointer(o.Obj())))
}

func Unref(o Object) {
	C.g_object_unref(C.gpointer(o.Obj()))
}

func SetProperty(o Object, name string, val interface{}) {
	n := (*C.gchar)(C.CString(name))
	defer C.free(unsafe.Pointer(n))
	C.g_object_set_property(o.Obj(), n, ValueOf(val).Val())
}

func Emit(o Object, sig Signal, args ...interface{}) interface{} {
	prms := make([]Value, len(args) + 1)
	prms[0] = *ValueOf(o.Obj())
	for i, a := range args {
		prms[i+1] = *ValueOf(a)
	}
	ret := new(Value)
	C._signal_emit(prms[0].Val(), C.guint(sig), ret.Val())
	fmt.Println("*** emitl ***")
	return ret.Get()
}

var handlers = map[uintptr]map[SigHandlerId]*reflect.Value{}

func Connect(o Object, sig Signal, cb_func interface{}) {
	cb := reflect.ValueOf(cb_func)
	if cb.Kind() != reflect.Func {
		panic("cb_func is not a function")
	}
	// Check that function parameters and return value match to signal
	var sq C.GSignalQuery
	C.g_signal_query(C.guint(sig), &sq)
	ft := cb.Type()
	if ft.NumOut() > 1 || ft.NumOut()==1 && Type(sq.return_type) == TYPE_NONE {
		panic("Number of function return values doesn't match signal spec.")
	}
	if ft.NumIn() != int(sq.n_params) {
		panic("Number of function parameters doesn't match signal spec.")
	}
	if ft.NumOut() != 0 && !Type(sq.return_type).Match(ft.Out(0)) {
		panic("Type of function return value doesn't match signal spec.")
	}
	pt := (*[1<<16]Type)(unsafe.Pointer(sq.param_types))[:int(sq.n_params)]
	for i := 0; i < ft.NumIn(); i++ {
		if !pt[i].Match(ft.In(i)) {
			panic(fmt.Sprintf(
				"%d function parameter type doesn't match signal spec.", i,
			))
		}
	}
	// Setup closure and connect it to signal
	cl := C._closure_new(o.Obj())
	cl.h_id = C._signal_connect(o.Obj(), C.guint(sig), cl)
	ptr := uintptr(unsafe.Pointer(o.Obj()))
	oh := handlers[ptr]
	if oh == nil {
		oh = map[SigHandlerId]*reflect.Value{}
		handlers[ptr] = oh
	}
	oh[SigHandlerId(cl.h_id)] = &cb
}

/*typedef struct {
	GClosure *cl;
	GValue *ret_val;
	guint n_param;
	const GValue *params;
	gpointer ih;
	gpointer data;
} MarshalParams;*/

//export go_marshal
func marshal(mp unsafe.Pointer) {
	p := (*C.MarshalParams)(mp)
	fmt.Println("*** marshal ***")
	cl := (*C.GoClosure)(p.cl)
	cb := handlers[uintptr(cl.o)][SigHandlerId(cl.h_id)]

	// TU_SKONCZYLEM
	//cb.Call(in)
	fmt.Println("cb", cb)
	fmt.Println("ret_val", p.ret_val)
}


// Binding for generic GObjec

type Params map[string]interface{}

type Obj C.GObject

// Implementation of Object interface

func (o *Obj) Obj() *C.GObject {
	return (*C.GObject)(o)
}

func (o *Obj) FromPtr(p C.gpointer) Object {
	return (*Obj)(o)
}

// Returns C pointer
func NewObj(t Type, params Params) *Obj {
	if params == nil || len(params) == 0 {
		return (*Obj)(C.g_object_newv(C.GType(t), 0, nil))
	}
	p := make([]C.GParameter, len(params))
	i := 0
	for k, v := range params {
		p[i].name = (*C.gchar)(C.CString(k))
		defer C.free(unsafe.Pointer(p[i].name))
		p[i].value = C.GValue(*ValueOf(v))
		i++
	}
	return (*Obj)(C.g_object_newv(C.GType(t), C.guint(i), &p[0]))
}
