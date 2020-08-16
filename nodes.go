package nohtml

import (
	"strings"

	"golang.org/x/net/html"
)

func isNav(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	return n.Data == "body" || (n.Data == "div" && hasClass(n, "nav-content"))
}

func isEmpty(n *html.Node) bool {
	return n.FirstChild == nil
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return contains(strings.Fields(attr.Val), class)
		}
	}
	return false
}

func addClass(n *html.Node, class string) *html.Node {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			attr.Val += " " + class
			return n
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: "class", Val: class})
	return n
}

func contains(sx []string, s string) bool {
	for _, s1 := range sx {
		if s1 == s {
			return true
		}
	}
	return false
}
