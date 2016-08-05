package headers

import (
	"net/http"

	"github.com/chrisolsen/quince"
	"golang.org/x/net/context"
)

// Set sets the response header to the key and value provided
func Set(key, value string) quince.Middleware {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		w.Header().Set(key, value)
		return c
	}
}
