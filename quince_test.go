package quince

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"

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

		return nil
	}

	c := context.Background()
	list := New(mw1, mw2)

	list.Run(c, nil, nil)
}

func Test_URLParams(t *testing.T) {
	router := httprouter.New()

	http.Handle("/", router)
}
