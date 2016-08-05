package basicauth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisolsen/quince"
	"golang.org/x/net/context"
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

	mw := Authenticate(func(c context.Context, user, password string) (context.Context, bool) {
		return c, false
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

	mw := Authenticate(func(c context.Context, user, password string) (context.Context, bool) {
		return c, true
	})
	mw(c, w, r)

	if w.Code != http.StatusOK {
		t.Error("Invalid response status: ", w.Code)
		return
	}
}

func Test_StoreAuthTokenInContext(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	c := appengine.NewContext(r)
	w := httptest.NewRecorder()

	r.Header.Add("Authorization", `Basic Zm9vOmJhcg==`)

	mw1 := Authenticate(func(c context.Context, user, password string) (context.Context, bool) {
		token := "some_token_obtained_within_app"
		c = context.WithValue(c, "token", token)
		return c, true
	})

	mw2 := func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		token := c.Value("token")
		if token == nil {
			t.Error("No token found in the context")
			return c
		}
		if token.(string) != "some_token_obtained_within_app" {
			t.Error("token does not match")
			return c
		}
		return c
	}

	q := quince.New(mw1, mw2)
	q.Run(c, w, r)
}
