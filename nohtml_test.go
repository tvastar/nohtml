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

func TestExamples(t *testing.T) {
	examples := map[string][2]interface{}{
		"simple string": {"hello", "<body>hello</body>"},
		"error":         {errors.New("foo"), "<body><pre>foo</pre></body>"},
	}

	for name, inOut := range examples {
		input, expected := inOut[0], inOut[1]
		t.Run(name, func(t *testing.T) {
			got, err := curl("GET", "/", "", input)
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
