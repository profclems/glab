package httpmock

import (
	"fmt"
	"net/http"
	"sync"
)

type matchType int

const (
	PathOnly matchType = iota
	HostOnly
	FullURL
	HostAndPath
	PathAndQuerystring
)

type Mocker struct {
	mu       sync.Mutex
	stubs    []*Stub
	Requests []*http.Request
	MatchURL matchType // if false, only matches the path and if true, matches full url
}

func New() *Mocker {
	return &Mocker{}
}

func (r *Mocker) RegisterResponder(method, path string, resp Responder) {
	r.stubs = append(r.stubs, &Stub{
		Matcher:   newRequest(method, path, r.MatchURL),
		Responder: resp,
	})
}

type Testing interface {
	Errorf(string, ...interface{})
	Helper()
}

func (r *Mocker) Verify(t Testing) {
	n := 0
	for _, s := range r.stubs {
		if !s.matched {
			n++
		}
	}
	if n > 0 {
		t.Helper()
		t.Errorf("%d unmatched HTTP stubs", n)
	}
}

// RoundTrip satisfies http.RoundTripper
func (r *Mocker) RoundTrip(req *http.Request) (*http.Response, error) {
	var stub *Stub

	r.mu.Lock()
	for _, s := range r.stubs {
		if s.matched || !s.Matcher(req) {
			continue
		}
		if stub != nil {
			r.mu.Unlock()
			return nil, fmt.Errorf("more than 1 stub matched %v", req)
		}
		stub = s
	}
	if stub != nil {
		stub.matched = true
	}

	if stub == nil {
		r.mu.Unlock()
		return nil, fmt.Errorf("no registered stubs matched %v", req)
	}

	r.Requests = append(r.Requests, req)
	r.mu.Unlock()

	return stub.Responder(req)
}
