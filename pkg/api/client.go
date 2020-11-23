package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/profclems/glab/internal/glinstance"
	"github.com/xanzy/go-gitlab"
)

var (
	apiClient *gitlab.Client
	err       error
	Protocol  = "https"
)

// Init initializes a gitlab client for use throughout glab.
func Init(host, token string, allowInsecure bool, isGraphQL bool) (*gitlab.Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: allowInsecure,
			},
		},
	}
	return gitlabClient(httpClient, token, host, isGraphQL)
}

func InitWithCustomCA(host, token, caFile string, isGraphQL bool) (*gitlab.Client, error) {
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	// use system cert pool as a baseline
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	caCertPool.AppendCertsFromPEM(caCert)

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return gitlabClient(httpClient, token, host, isGraphQL)
}

func gitlabClient(httpClient *http.Client, token, host string, isGraphQL bool) (*gitlab.Client, error) {
	var baseURL string
	if host == "" {
		host = glinstance.OverridableDefault()
	}
	if isGraphQL {
		baseURL = glinstance.GraphQLEndpoint(host, Protocol)
	} else {
		baseURL = glinstance.APIEndpoint(host, Protocol)
	}
	apiClient, err = gitlab.NewClient(token, gitlab.WithHTTPClient(httpClient), gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GitLab client: %v", err)
	}
	return apiClient, nil
}

func TestClient(httpClient *http.Client, token, host string, isGraphQL bool) (*gitlab.Client, error) {
	return gitlabClient(httpClient, token, host, isGraphQL)
}
