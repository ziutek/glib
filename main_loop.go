package glib

/*
#include <glib-object.h>
*/
import "C"

type MainLoop struct {
	Object
}

func (l MainLoop) GMainLoop() *C.GMainLoop {
	return (*C.GMainLoop)(l.GPointer())
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
	k.Set(Pointer(C.g_main_loop_get_context(l.GMainLoop())))
	return k
}

func NewMainLoop(ctx *MainContext) *MainLoop {
	l := new(MainLoop)
	l.Set(Pointer(C.g_main_loop_new(ctx.GMainContext(), 0)))
	return l
}
