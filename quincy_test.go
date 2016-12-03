package quincy

import (
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

func Test_Context(t *testing.T) {
	mw1 := func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		return context.WithValue(c, "key", "foobar")
	}

	mw2 := func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		val := c.Value("key")
		if val == nil {
			t.Error("no value set")
			return nil
		}

		if val.(string) != "foobar" {
			t.Error("value does not match the expected")
			return nil
		}

		return c
	}

	c := context.Background()
	list := New(mw1, mw2)

	list.Run(c, nil, nil)
}

func Test_MiddlewareCancel(t *testing.T) {

	mw1 := func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		c, cancel := context.WithCancel(c)
		cancel()
		return c
	}

	mw2 := func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		t.Error("Should not make it to this middleware")
		return c
	}

	c := context.Background()
	list := New(mw1, mw2)

	list.Run(c, nil, nil)
}

func Test_NoMiddleware(t *testing.T) {
	c := context.Background()
	q := New()
	q.Then(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		// nothing to see here, this test passes by not blowing up
	})
	q.Run(c, nil, nil)
}
