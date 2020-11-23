package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/xanzy/go-gitlab"
)

type HTTPResponse struct {
	Response *gitlab.Response
	Output   *bytes.Buffer
	Request  *retryablehttp.Request
}

func httpRequest(client *gitlab.Client, config config.Config, hostname string, method string, p string, params interface{}, headers []string) (*HTTPResponse, error) {
	var err error
	isGraphQL := p == "graphql"
	if client.BaseURL().Host != hostname || isGraphQL {
		client, err = cmdutils.HttpClientFunc(hostname, config, isGraphQL)
		if err != nil {
			return nil, err
		}
	}

	baseURL := client.BaseURL()
	baseURLStr := baseURL.String()
	if strings.Contains(p, "://") {
		baseURLStr = p
	} else if isGraphQL {
		baseURL.Path = strings.Replace(baseURL.Path, "api/v4/", "", 1)
		baseURLStr = baseURL.String()
	} else {
		baseURLStr = baseURLStr + strings.TrimPrefix(p, "/")
	}
	var body io.Reader
	var bodyIsJSON bool
	switch pp := params.(type) {
	case map[string]interface{}:
		if strings.EqualFold(method, "GET") || strings.EqualFold(method, "DELETE") {
			baseURLStr, err = parseQuery(baseURLStr, pp)
			if err != nil {
				return nil, err
			}
		} else {
			for key, value := range pp {
				switch vv := value.(type) {
				case []byte:
					pp[key] = string(vv)
				}
			}
			if isGraphQL {
				pp = groupGraphQLVariables(pp)
			}
			b, err := json.Marshal(pp)
			if err != nil {
				return nil, fmt.Errorf("error serializing parameters: %w", err)
			}
			body = bytes.NewBuffer(b)
			bodyIsJSON = true
		}
	case io.Reader:
		body = pp
	case nil:
		body = nil
	default:
		return nil, fmt.Errorf("unrecognized parameters type: %v", params)
	}

	baseURL, _ = url.Parse(baseURLStr)
	req, err := newRequest(method, baseURL, body, client.UserAgent, headers, bodyIsJSON)

	if err != nil {
		return nil, err
	}

	hr := &HTTPResponse{
		Output:   &bytes.Buffer{},
		Response: &gitlab.Response{},
	}
	resp, err := client.Do(req, hr.Output)
	if err != nil {
		return nil, err
	}
	hr.Response = resp
	hr.Request = req
	return hr, err
}

func groupGraphQLVariables(params map[string]interface{}) map[string]interface{} {
	topLevel := make(map[string]interface{})
	variables := make(map[string]interface{})

	for key, val := range params {
		switch key {
		case "query", "operationName":
			topLevel[key] = val
		default:
			variables[key] = val
		}
	}

	if len(variables) > 0 {
		topLevel["variables"] = variables
	}
	return topLevel
}

func parseQuery(path string, params map[string]interface{}) (string, error) {
	if len(params) == 0 {
		return path, nil
	}
	q := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case string:
			q.Add(key, v)
		case []byte:
			q.Add(key, string(v))
		case nil:
			q.Add(key, "")
		case int:
			q.Add(key, fmt.Sprintf("%d", v))
		case bool:
			q.Add(key, fmt.Sprintf("%v", v))
		default:
			return "", fmt.Errorf("unknown type %v", v)
		}
	}

	sep := "?"
	if strings.ContainsRune(path, '?') {
		sep = "&"
	}
	return path + sep + q.Encode(), nil
}

func newRequest(method string, baseURL *url.URL, body io.Reader, userAgent string, headers []string, bodyIsJSON bool) (*retryablehttp.Request, error) {
	req, err := http.NewRequest(method, baseURL.String(), body)
	if err != nil {
		return nil, err
	}

	for _, h := range headers {
		idx := strings.IndexRune(h, ':')
		if idx == -1 {
			return nil, fmt.Errorf("header %q requires a value separated by ':'", h)
		}
		name, value := h[0:idx], strings.TrimSpace(h[idx+1:])
		if strings.EqualFold(name, "Content-Length") {
			length, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return nil, err
			}
			req.ContentLength = length
		} else {
			req.Header.Add(name, value)
		}
	}
	if bodyIsJSON && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	//req, err := retryablehttp.NewRequest(method, baseURL.String(), body)
	rReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	return rReq, nil
}
