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

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
)

func httpRequest(client *api.Client, config config.Config, hostname string, method string, p string, params interface{}, headers []string) (*http.Response, error) {
	var err error
	isGraphQL := p == "graphql"
	if client.Lab().BaseURL().Host != hostname || isGraphQL {
		fmt.Println("tes")
		client, err = api.NewClientWithCfg(hostname, config, isGraphQL)
		if err != nil {
			return nil, err
		}
	}

	baseURL := client.Lab().BaseURL()
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
	req, err := newRequest(client, method, baseURL, body, client.LabClient.UserAgent, headers, bodyIsJSON)

	if err != nil {
		return nil, err
	}
	return client.HTTPClient().Do(req)
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

func newRequest(c *api.Client, method string, baseURL *url.URL, body io.Reader, userAgent string, headers []string, bodyIsJSON bool) (*http.Request, error) {
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

	// TODO: support GITLAB_CI_TOKEN
	switch c.AuthType {
	case api.OAuthToken:
		req.Header.Set("Authorization", "Bearer "+c.Token())
	case api.PrivateToken:
		req.Header.Set("PRIVATE-TOKEN", c.Token())
	}

	return req, nil
}
