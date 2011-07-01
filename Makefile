include $(GOROOT)/src/Make.inc

TARG = glib
CGOFILES = pkgconfig.go type.go value.go object.go signal.go

include $(GOROOT)/src/Make.pkg
