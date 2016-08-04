package urlparams

import (
	"errors"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

const (
	paramsKey = "urlparams"
)

// Get obtains the httprouter.Params struct
func Get(c context.Context) (httprouter.Params, error) {
	p := c.Value(paramsKey)
	if p == nil {
		return nil, errors.New("No params were available within the context")
	}
	return p.(httprouter.Params), nil
}

// ByName obtains the url params previously set within the context. This is a helper
// func that should only be used if only obtaining one value.
func ByName(c context.Context, key string) (string, error) {
	p := c.Value(paramsKey)
	if p == nil {
		return "", errors.New("No params were available within the context")
	}
	return p.(httprouter.Params).ByName(key), nil
}

// Put injects the url params into the context
func Put(c context.Context, params httprouter.Params) context.Context {
	return context.WithValue(c, paramsKey, params)
}
