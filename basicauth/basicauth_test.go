package basicauth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func Test_InitialRequestWithMissingHeaders(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	c := appengine.NewContext(r)
	w := httptest.NewRecorder()

	mw := Authenticate(nil)
	mw(c, w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("Invalid response status: ", w.Code)
		return
	}

	header := w.Header().Get("www-authenticate")
	if len(header) == 0 {
		t.Error("missing www-authenticate header")
		return
	}

	if strings.Index(header, "Basic") != 0 {
		t.Error("Missing `Basic` in header")
		return
	}
}

func Test_InvalidCredentials(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	c := appengine.NewContext(r)
	w := httptest.NewRecorder()

	r.Header.Add("Authorization", `Basic Zm9vOmJhcg==`)

	mw := Authenticate(func(user, password string) bool {
		return false
	})
	mw(c, w, r)

	if w.Code != http.StatusUnauthorized {
		t.Error("Invalid response status: ", w.Code)
		return
	}
}

func Test_ValidCredentials(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	c := appengine.NewContext(r)
	w := httptest.NewRecorder()

	r.Header.Add("Authorization", `Basic Zm9vOmJhcg==`)

	mw := Authenticate(func(user, password string) bool {
		return true
	})
	mw(c, w, r)

	if w.Code != http.StatusOK {
		t.Error("Invalid response status: ", w.Code)
		return
	}
}
