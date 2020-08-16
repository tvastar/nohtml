package nohtml

import (
	"context"
	"net/http"
	"strings"
)

// New creates a context from a http request.
func New(r *http.Request) *Context {
	p := strings.Split(r.URL.Path, "/")
	h := strings.Split(r.URL.Fragment, "/")
	return &Context{r, nil, p, nil, h}
}

// Context tracks a request context.
type Context struct {
	r *http.Request

	path, rpath []string
	hash, rhash []string
}

// Context returns the context.Context associated with the current request.
func (c *Context) Context() context.Context {
	return c.r.Context()
}

// Path returns the first element in the request path and
// returns a sub-context (which has the rest of the path).
func (c *Context) Path() (string, *Context) {
	for idx, p := range c.rpath {
		if p != "" {
			result := c.AppendPath(p)
			result.rpath = c.rpath[idx+1:]
			return p, result
		}
	}
	return "", c
}

// AppendPath returns a new context with the provided path element
// appended to it.
func (c *Context) AppendPath(pathElement string) *Context {
	result := *c
	result.path = append(append([]string{}, result.path...), pathElement)
	return &result
}

// Hash returns the first element in the request hash and
// returns a sub-context (which has the rest of the hash).
func (c *Context) Hash() (string, *Context) {
	for idx, h := range c.rhash {
		if h != "" {
			result := c.AppendHash(h)
			result.rhash = c.rhash[idx+1:]
			return h, result
		}
	}
	return "", c
}

// AppendHash returns a new context with the provided hash element
// appended to it.
func (c *Context) AppendHash(hashElement string) *Context {
	result := *c
	result.path = append(append([]string{}, result.hash...), hashElement)
	return &result
}
