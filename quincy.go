package quincy

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Middleware is a http.HandlerFunc that also includes a context and url params variables
type Middleware func(context.Context, http.ResponseWriter, *http.Request) context.Context

// HandlerFunc ...
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// Q allows a list middleware functions to be created and run
type Q struct {
	fns []Middleware
}

// New initializes the middleware chain with one or more handler functions.
// The returned pointer allows for additional middleware methods to be added or for the chain to be run.
//	q := que.New(foo, bar)
func New(fns ...Middleware) *Q {
	q := Q{}
	q.fns = fns
	return &q
}

// Add allows for one or more middleware handler functions to be added to the existing chain
//	q := que.New(cors, format)
//	q.Add(auth)
func (q *Q) Add(fns ...Middleware) {
	q.fns = append(q.fns, fns...)
}

// Run executes the handler chain, which is most useful in tests
//	q := que.New(foo, bar)
// 	q.Add(func(c context.Context, w http.ResponseWriter, r *http.Request) {
// 		// perform tests here
// 	})
//  inst := aetest.NewInstance(nil)
// 	r := inst.NewRequest("GET", "/", nil)
// 	w := httpTest.NewRecorder()
// 	c := appengine.NewContext(r)
// 	q.Run(c, w, r)
func (q *Q) Run(c context.Context, w http.ResponseWriter, r *http.Request) {
	chain(q.fns)(c, w, r)
}

// Then returns the chain of existing middleware that includes the final HandlerFunc argument.
//	q := que.New(foo, bar)
//  router.Get("/", q.Then(handleRoot))
func (q *Q) Then(fn HandlerFunc) func(http.ResponseWriter, *http.Request) {
	chn := chain(q.fns)

	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		c = chn(c, w, r)

		if c.Err() == nil {
			fn(c, w, r)
		}
	}
}

func chain(fns []Middleware) Middleware {
	var next Middleware
	var count = len(fns)
	for i := count - 1; i >= 0; i-- {
		next = link(fns[i], next)
	}

	return next
}

func link(current, next Middleware) Middleware {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		c = current(c, w, r)
		if c.Err() != nil {
			return c
		}
		if next != nil {
			c = next(c, w, r)
		}
		return c
	}
}
