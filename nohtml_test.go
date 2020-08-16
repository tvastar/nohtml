package nohtml_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/tvastar/nohtml"
)

func TestGetSuccesses(t *testing.T) {
	// testName => [testInput, testOutput]
	examples := map[string][3]interface{}{
		"simple string": {"", "hello", "<body>hello</body>"},
		"error":         {"", errors.New("foo"), "<body><pre>foo</pre></body>"},
		"simple map": {
			"",
			map[string]interface{}{"hello": 42, "world": 22},
			`<body><div class="nav-options"><a href="/hello">hello</a><a href="/world">world</a></div><div class="nav-content"></div></body>`,
		},
		"simple map, path": {
			"hello",
			map[string]interface{}{"hello": 42, "world": 22},
			`<body><div class="nav-options"><a href="/hello" class="selected">hello</a><a href="/world">world</a></div><div class="nav-content">42</div></body>`,
		},
		"nested path": {
			"hello/world",
			map[string]interface{}{"hello": map[string]string{"world": "foo"}},
			`<body><div class="nav-options"><a href="/hello" class="selected">hello</a></div><div class="nav-content"><div class="nav-options"><a href="/hello/world" class="selected">world</a></div><div class="nav-content">foo</div></div></body>`,
		},
	}

	for name, inOut := range examples {
		path, input, expected := inOut[0], inOut[1], inOut[2]
		t.Run(name, func(t *testing.T) {
			got, err := curl("GET", path.(string), "", input)
			if err != nil {
				t.Error("Unexpected", err)
			}
			if got != expected {
				t.Error("Got", got, "Expected", expected)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	// testName => [testInput, method, path, body, testOutput]
	examples := map[string][5]interface{}{
		"invalid path: string": {"hello", "GET", "/foo", "", "<body><pre>unknown path: foo</pre></body>"},
		"post: string":         {"hello", "POST", "/", "", "<body><pre>only http GET supported</pre></body>"},
	}

	for name, inOut := range examples {
		input, method, path, body, expected := inOut[0], inOut[1], inOut[2], inOut[3], inOut[4]
		t.Run(name, func(t *testing.T) {
			got, err := curl(method.(string), path.(string), body.(string), input)
			if err != nil {
				t.Error("Unexpected", err)
			}
			if got != expected {
				t.Error("Got", got, "Expected", expected)
			}
		})
	}
}

func curl(method, endpoint, body string, v interface{}) (string, error) {
	ts := httptest.NewServer(nohtml.Handler(v))
	defer ts.Close()

	r := strings.NewReader(body)
	req, err := http.NewRequest(method, ts.URL+"/"+endpoint, r)
	if err != nil {
		return "", err
	}

	if body != "" {
		req.Header.Add("Content-Type", "application/x-www-url-encoded")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return string(data), err
}
