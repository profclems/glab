package httpmock

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Matcher func(req *http.Request) bool
type Responder func(req *http.Request) (*http.Response, error)

type Stub struct {
	matched   bool
	Matcher   Matcher
	Responder Responder
}

func MatchAny(*http.Request) bool {
	return true
}

func newRequest(method, path string, match matchType) Matcher {
	return func(req *http.Request) bool {
		if !strings.EqualFold(req.Method, method) {
			return false
		}
		if match == PathOnly {
			if !strings.HasPrefix(path, "/api/v4") {
				path = "/api/v4" + path
			}
			return req.URL.Path == path
		}
		u, err := url.Parse(path)
		if err != nil {
			return false
		}
		if match == FullURL {
			return req.URL.String() == u.String()
		}
		if match == HostOnly {
			return req.URL.Host == u.Host
		}
		if match == HostAndPath {
			return req.URL.Host == u.Host && req.URL.Path == u.Path
		}
		if match == PathAndQuerystring {
			return req.URL.RawQuery == u.RawQuery && req.URL.Path == u.Path
		}
		return false
	}
}

func NewStringResponse(status int, body string) Responder {
	return func(req *http.Request) (*http.Response, error) {
		return httpResponse(status, req, bytes.NewBufferString(body)), nil
	}
}

func NewJSONResponse(status int, body interface{}) Responder {
	return func(req *http.Request) (*http.Response, error) {
		b, _ := json.Marshal(body)
		return httpResponse(status, req, bytes.NewBuffer(b)), nil
	}
}

func NewFileResponse(status int, filename string) Responder {
	return func(req *http.Request) (*http.Response, error) {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		return httpResponse(status, req, f), nil
	}
}

func httpResponse(status int, req *http.Request, body io.Reader) *http.Response {
	return &http.Response{
		StatusCode: status,
		Request:    req,
		Body:       ioutil.NopCloser(body),
	}
}
