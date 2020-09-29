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
)

// Init initializes a gitlab client for use throughout glab.
func Init(host, token string, allowInsecure bool) (*gitlab.Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
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

	apiClient, err = gitlab.NewClient(token, gitlab.WithHTTPClient(httpClient), gitlab.WithBaseURL(glinstance.APIEndpoint(host)))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GitLab client: %v", err)
	}
	return apiClient, nil
}

// InitWithBasicAuth initialises a client with username and password.
func InitWithBasicAuth(host, username, password string) (*gitlab.Client, error) {
	apiClient, err = gitlab.NewBasicAuthClient(
		username,
		password,
		gitlab.WithBaseURL(glinstance.APIEndpoint(host)),
	)
	if err != nil {
		return nil, err
	}
	return apiClient, nil
}

func InitWithCustomCA(host, token, caFile string) (*gitlab.Client, error) {
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
				DualStack: true,
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

	apiClient, _ = gitlab.NewClient(token, gitlab.WithHTTPClient(httpClient), gitlab.WithBaseURL(glinstance.APIEndpoint(host)))
	return apiClient, nil
}
