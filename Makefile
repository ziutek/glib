include $(GOROOT)/src/Make.inc

TARG = github.com/ziutek/glib
CGOFILES = type.go value.go object.go signal.go main_context.go main_loop.go

include $(GOROOT)/src/Make.pkg
