#include <glib-object.h>

typedef struct {
	GClosure cl;
	gpointer cb;
} GoClosure;
