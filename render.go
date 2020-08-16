package nohtml

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/pkg/errors"

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

	if c.r.Method != http.MethodGet {
		c.Render(errors.New("only http GET supported"), parent)
		return
	}

	k := v.Kind()
	switch {
	case k == reflect.Interface:
		c.render(v.Elem(), parent)
	case k == reflect.Map && isNav(parent) && isEmpty(parent):
		c.renderNav(v, parent)

	default:
		path, _ := c.Path()
		if path != "" {
			c.Render(errors.New("unknown path: "+path), parent)
			return
		}

		c.renderText(v, parent)
	}
}

func (c *Context) renderText(v reflect.Value, parent *html.Node) {
	parent.AppendChild(c.Text(fmt.Sprintf("%v", v.Interface())))
}

func (c *Context) renderNav(v reflect.Value, parent *html.Node) {
	options, contents := c.RenderNavContainer(parent)
	path, cx := c.Path()
	strkeys := []string{}
	keys := map[string]reflect.Value{}
	for _, key := range v.MapKeys() {
		s := fmt.Sprintf("%v", key.Interface())
		strkeys = append(strkeys, s)
		keys[s] = key
	}
	sort.Strings(strkeys)

	for _, s := range strkeys {
		key := keys[s]
		option := c.Link(s, c.HRef()+"/"+s)
		if s == path {
			option = addClass(option, "selected")
			cx.render(v.MapIndex(key), contents)
		}
		options.AppendChild(option)
	}
}

// RenderNavContainer renders a simple nav container with a sub
// container for options and contents.
func (c *Context) RenderNavContainer(parent *html.Node) (options, contents *html.Node) {
	options = addClass(c.Element("div"), "nav-options")
	contents = addClass(c.Element("div"), "nav-content")
	parent.AppendChild(options)
	parent.AppendChild(contents)
	return options, contents
}

// RenderLink returns an anchor node.
func (c *Context) Link(label, href string) *html.Node {
	l := c.Element("a", c.Text(label))
	l.Attr = append(l.Attr, html.Attribute{Key: "href", Val: href})
	return l
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
