package nohtml

import (
	"fmt"
	"reflect"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Render renders any value, appending into the provided HTML node.
func (c *Context) Render(v interface{}, parent *html.Node) {
	c.render(reflect.ValueOf(v), parent)
}

func (c *Context) render(v reflect.Value, parent *html.Node) {
	switch v := v.Interface().(type) {
	case error:
		parent.AppendChild(c.Element("pre", c.Text(fmt.Sprintf("%v", v))))
		return
	}

	parent.AppendChild(c.Text(fmt.Sprintf("%v", v.Interface())))
}

// Element returns a html Element node with the provided tag and children.
func (c *Context) Element(tag string, children ...*html.Node) *html.Node {
	a := atom.Lookup([]byte(tag))
	n := &html.Node{Type: html.ElementNode, Data: tag, DataAtom: a}
	for _, c := range children {
		n.AppendChild(c)
	}
	return n
}

// Text returns a html Text node with the provided content.
func (c *Context) Text(content string) *html.Node {
	return &html.Node{Type: html.TextNode, Data: content}
}
