// Package nohtml provides functions to build web apps without HTML.
package nohtml

import (
	"net/http"

	"golang.org/x/net/html"
)

// Handler returns a http.Handler which serves the provided object.
func Handler(v interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := New(r)
		body := c.Element("body")
		c.Render(v, body)
		if err := html.Render(w, body); err != nil {
			panic(err)
		}
	})
}
