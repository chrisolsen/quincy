# quincy 

# DEPRECATED: This repo has been merged into github.com/chrisolsen/ae

quincy allows middleware sequenced functionality for AppEngine Go services. Url params are made available, via the `urlparams` package, by using the Go context library that is already used within AppEngine.

Since `context` is not available via the `http.Request` until Go 1.7, and AppEngine usually lags behind in regards to the version of Go made available, this library fills the hole until then.

## Example

```Go
import (
    "github.com/chrisolsen/quincy"
    "github.com/chrisolsen/quincy/urlparams"
)

func init() {
    r := httprouter.New()
    q := quincy.New(middleware1, middleware2)

    r.GET("/accounts/{id}", q.Then(accountHandler))
}

func middleware1(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
    // access url params with helper method
    id := urlparams.ByName("id") 

    // get context values
    foo := c.ValueOf("foo") 

    // set context param
    c = context.WithValue(c, "key", "some value")

    return c
}

func accountHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
    // access context values here the same as the middleware    
}
```

## License

MIT License
