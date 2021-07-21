package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/xanzy/go-gitlab"
)

// AuthType represents an authentication type within GitLab.
type authType int

const (
	NoToken authType = iota
	OAuthToken
	PrivateToken
)

const UserAgent = "GLab - GitLab CLI"

// Global api client to be used throughout glab
var apiClient *Client

// Client represents an argument to NewClient
type Client struct {
	// LabClient represents GitLab API client.
	// Note: this is exported for tests. Do not access it directly. Use Lab() method
	LabClient *gitlab.Client
	// internal http client
	httpClient *http.Client
	// internal http client overrider
	httpClientOverride *http.Client
	// Token type used to make authenticated API calls.
	AuthType authType
	// custom certificate
	caFile string
	// Protocol: host url protocol to make requests. Default is https
	Protocol string

	host  string
	token string

	isGraphQL          bool
	allowInsecure      bool
	refreshLabInstance bool
}

func init() {
	// initialise the global api client to be used throughout glab
	RefreshClient()
}

// RefreshClient re-initializes the api client
func RefreshClient() {
	apiClient = &Client{
		Protocol:           "https",
		AuthType:           NoToken,
		httpClient:         &http.Client{},
		refreshLabInstance: true,
	}
}

// GetAPIClient returns the global DotEnv instance.
func GetClient() *Client {
	return apiClient
}

// HTTPClient returns the httpClient instance used to initialise the gitlab api client
func HTTPClient() *http.Client { return apiClient.HTTPClient() }
func (c *Client) HTTPClient() *http.Client {
	if c.httpClientOverride != nil {
		return c.httpClientOverride
	}
	if c.httpClient != nil {
		return c.httpClient
	}
	return &http.Client{}
}

// OverrideHTTPClient overrides the default http client
func OverrideHTTPClient(client *http.Client) { apiClient.OverrideHTTPClient(client) }
func (c *Client) OverrideHTTPClient(client *http.Client) {
	c.httpClientOverride = client
}

// Token returns the authentication token
func Token() string { return apiClient.Token() }
func (c *Client) Token() string {
	return c.token
}

func SetProtocol(protocol string) { apiClient.SetProtocol(protocol) }
func (c *Client) SetProtocol(protocol string) {
	c.Protocol = protocol
}

// NewClient initializes a api client for use throughout glab.
func NewClient(host, token string, allowInsecure bool, isGraphQL bool) (*Client, error) {
	apiClient.host = host
	apiClient.token = token
	apiClient.allowInsecure = allowInsecure
	apiClient.isGraphQL = isGraphQL

	if apiClient.httpClientOverride == nil {
		apiClient.httpClient = &http.Client{
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
					InsecureSkipVerify: apiClient.allowInsecure,
				},
			},
		}
	}
	apiClient.refreshLabInstance = true
	err := apiClient.NewLab()
	return apiClient, err
}

// NewClientWithCustomCA initializes the global api client with a self-signed certificate
func NewClientWithCustomCA(host, token, caFile string, isGraphQL bool) (*Client, error) {
	apiClient.host = host
	apiClient.token = token
	apiClient.caFile = caFile
	apiClient.isGraphQL = isGraphQL

	if apiClient.httpClientOverride == nil {
		caCert, err := ioutil.ReadFile(apiClient.caFile)
		if err != nil {
			return nil, fmt.Errorf("error reading cert file: %w", err)
		}
		// use system cert pool as a baseline
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)

		apiClient.httpClient = &http.Client{
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
	}
	apiClient.refreshLabInstance = true
	err := apiClient.NewLab()
	return apiClient, err
}

// NewClientWithCustomCAClientCert initializes the global api client with a self-signed certificate and client certificates
func NewClientWithCustomCAClientCert(host, token, caFile string, certFile string, keyFile string, isGraphQL bool) (*Client, error) {
	apiClient.host = host
	apiClient.token = token
	apiClient.caFile = caFile
	apiClient.isGraphQL = isGraphQL

	if apiClient.httpClientOverride == nil {
		caCert, err := ioutil.ReadFile(apiClient.caFile)
		if err != nil {
			return nil, fmt.Errorf("error reading cert file: %w", err)
		}
		// use system cert pool as a baseline
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)

		clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		clientCerts := []tls.Certificate{clientCert}

		apiClient.httpClient = &http.Client{
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
					RootCAs:      caCertPool,
					Certificates: clientCerts,
				},
			},
		}
	}
	apiClient.refreshLabInstance = true
	err := apiClient.NewLab()
	return apiClient, err
}

// NewClientWithCfg initializes the global api with the config data
func NewClientWithCfg(repoHost string, cfg config.Config, isGraphQL bool) (client *Client, err error) {
	if repoHost == "" {
		repoHost = glinstance.OverridableDefault()
	}

	apiHost, _ := cfg.Get(repoHost, "api_host")
	if apiHost == "" {
		apiHost = repoHost
	}

	apiProtocol, _ := cfg.Get(repoHost, "api_protocol")
	if apiProtocol != "" {
		SetProtocol(apiProtocol)
	}

	token, _ := cfg.Get(repoHost, "token")
	tlsVerify, _ := cfg.Get(repoHost, "skip_tls_verify")
	skipTlsVerify := tlsVerify == "true" || tlsVerify == "1"
	caCert, _ := cfg.Get(repoHost, "ca_cert")
	clientCert, _ := cfg.Get(repoHost, "client_cert")
	keyFile, _ := cfg.Get(repoHost, "client_key")
	if caCert != "" && clientCert != "" && keyFile != "" {
		client, err = NewClientWithCustomCAClientCert(apiHost, token, caCert, clientCert, keyFile, isGraphQL)
	} else if caCert != "" {
		client, err = NewClientWithCustomCA(apiHost, token, caCert, isGraphQL)
	} else {
		client, err = NewClient(apiHost, token, skipTlsVerify, isGraphQL)
	}
	return
}

// NewLab initializes the GitLab Client
func (c *Client) NewLab() error {
	var err error
	var baseURL string
	httpClient := c.httpClient

	if c.httpClientOverride != nil {
		httpClient = c.httpClientOverride
	}
	if apiClient.refreshLabInstance {
		if c.host == "" {
			c.host = glinstance.OverridableDefault()
		}
		if c.isGraphQL {
			baseURL = glinstance.GraphQLEndpoint(c.host, c.Protocol)
		} else {
			baseURL = glinstance.APIEndpoint(c.host, c.Protocol)
		}
		c.LabClient, err = gitlab.NewClient(c.token, gitlab.WithHTTPClient(httpClient), gitlab.WithBaseURL(baseURL))
		if err != nil {
			return fmt.Errorf("failed to initialize GitLab client: %v", err)
		}
		c.LabClient.UserAgent = UserAgent

		if c.token != "" {
			c.AuthType = PrivateToken
		}
	}
	return nil
}

// Lab returns the initialized GitLab client.
// Initializes a new GitLab Client if not initialized but error is ignored
func (c *Client) Lab() *gitlab.Client {
	if c.LabClient != nil {
		return c.LabClient
	}
	err := c.NewLab()
	if err != nil {
		c.LabClient = &gitlab.Client{}
	}
	return c.LabClient
}

// BaseURL returns a copy of the BaseURL
func (c *Client) BaseURL() *url.URL {
	return c.Lab().BaseURL()
}

func NewHTTPRequest(c *Client, method string, baseURL *url.URL, body io.Reader, headers []string, bodyIsJSON bool) (*http.Request, error) {
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

	if c.Lab().UserAgent != "" {
		req.Header.Set("User-Agent", c.Lab().UserAgent)
	}

	// TODO: support GITLAB_CI_TOKEN
	switch c.AuthType {
	case OAuthToken:
		req.Header.Set("Authorization", "Bearer "+c.Token())
	case PrivateToken:
		req.Header.Set("PRIVATE-TOKEN", c.Token())
	}

	return req, nil
}

func TestClient(httpClient *http.Client, token, host string, isGraphQL bool) (*Client, error) {
	testClient, err := NewClient(host, token, true, isGraphQL)
	if err != nil {
		return nil, err
	}
	testClient.SetProtocol("https")
	testClient.OverrideHTTPClient(httpClient)
	testClient.refreshLabInstance = true
	if token != "" {
		testClient.AuthType = PrivateToken
	}
	return testClient, nil
}
