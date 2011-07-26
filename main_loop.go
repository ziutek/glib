package glib

/*
#include <glib-object.h>
*/
import "C"

type MainLoop struct {
	Object
}

func (l MainLoop) GMainLoop() *C.GMainLoop {
	return (*C.GMainLoop)(l.GetPtr())
}

func (l MainLoop) Run() {
	C.g_main_loop_run(l.GMainLoop())
}

func (l MainLoop) Quit() {
	C.g_main_loop_quit(l.GMainLoop())
}

func (l MainLoop) IsRunning() bool {
	return C.g_main_loop_is_running(l.GMainLoop()) != 0
}

func (l MainLoop) GetContext() *MainContext {
	k := new(MainContext)
	k.SetPtr(Pointer(C.g_main_loop_get_context(l.GMainLoop())))
	return k
}

func NewMainLoop(ctx *MainContext) *MainLoop {
	l := new(MainLoop)
	var c *C.GMainContext
	if ctx != nil {
		c = ctx.GMainContext()
	}
	l.SetPtr(Pointer(C.g_main_loop_new(c, 0)))
	return l
}
