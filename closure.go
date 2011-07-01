package glib

/*
#include "closure.h"

static inline GoClosure* go_closure_ref(GoClosure* c) {
	return (GoClosure*) g_closure_ref((GClosure*) c);
}

static inline void go_closure_unref(GoClosure* c) {
	g_closure_unref((GClosure*) c);
}

typedef struct {
	GClosure *cl;
	GValue *ret_val;
	guint n_param;
	const GValue *params;
	gpointer ih;
	gpointer data;
} MarshalParams;

extern void go_marshal(gpointer mp);  

static void closure_marshal(GClosure* cl, GValue* ret_val, guint n_param,
		const GValue* params, gpointer ih, gpointer data) {
	MarshalParams mp = {cl, ret_val, n_param, params, ih, data};
	go_marshal(&mp);	
}

static GoClosure* go_closure_new(void* cb) {
	GoClosure *cl = (GoClosure*) g_closure_new_simple(sizeof (GoClosure), NULL);
	cl->cb = cb;
	g_closure_set_marshal((GClosure *) cl, closure_marshal);
	return cl;
}
*/
import "C"

import (
	"runtime"
	"reflect"
	"unsafe"
	"fmt"
)

type Closure struct {
	cl *C.GoClosure
	cb *reflect.Value // Callback function
}

func NewClosure(cb_func interface{}) (c *Closure) {
	cb := reflect.ValueOf(cb_func)
	if cb.Kind() != reflect.Func {
		panic("cb_func is not a function")
	}
	c = &Closure{C.go_closure_new(unsafe.Pointer(&cb)), &cb}
	runtime.SetFinalizer(c, (*Closure).destroy)
	return
}

func (c *Closure) Ref() *Closure {
	if c.cl == nil {
		panic("Ref on nil closure")
	}
	return &Closure{C.go_closure_ref(c.cl), c.cb}
}

func (c *Closure) Unref() {
	if c.cl == nil {
		panic("Unref on nil closure")
	}
	C.go_closure_unref(c.cl)
	c.cl = nil
	c.cb = nil
}

func (c *Closure) destroy() {
	if c.cl != nil {
		C.go_closure_unref(c.cl)
	}
}

//export go_marshal
func marshal(mp unsafe.Pointer) {
	//p := (*C.MarshalParams)(mp)
	fmt.Println("Doszedl!")
}
