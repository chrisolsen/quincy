package headers

import (
	"net/http/httptest"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func Test_SetHeaders(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	c := appengine.NewContext(r)
	w := httptest.NewRecorder()

	mw := Set("foo", "bar")
	mw(c, w, r)

	if w.Header().Get("foo") != "bar" {
		t.Error("Header not set")
		return
	}
}
