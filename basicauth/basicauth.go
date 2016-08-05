package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/chrisolsen/quince"
	"golang.org/x/net/context"
)

// Auth is the function type of the custom function that is required to perform the custom authentication
type Auth func(name, password string) bool

// Authenticate performs the authentication using the passed in auth function
func Authenticate(auth Auth) quince.Middleware {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		reject := func() context.Context {
			w.Header().Set("WWW-Authenticate", `Basic realm=""`)
			w.WriteHeader(http.StatusUnauthorized)
			c, cancel := context.WithCancel(c)
			defer cancel()
			return c
		}

		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			return reject()
		}

		raw := strings.TrimRight(authHeader[len("basic="):], "=")
		input, err := base64.RawStdEncoding.DecodeString(raw)
		if err != nil {
			return reject()
		}

		parts := strings.Split(string(input), ":")
		if len(parts) != 2 {
			return reject()
		}
		if !auth(parts[0], parts[1]) {
			return reject()
		}

		return c
	}
}
