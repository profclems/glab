package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/internal/config"
)

func httpRequest(client *api.Client, config config.Config, hostname string, method string, p string, params interface{}, headers []string) (*http.Response, error) {
	var err error
	isGraphQL := p == "graphql"
	if client.Lab().BaseURL().Host != hostname || isGraphQL {
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
				if vv, ok := value.([]byte); ok {
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
	req, err := api.NewHTTPRequest(client, method, baseURL, body, headers, bodyIsJSON)

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
