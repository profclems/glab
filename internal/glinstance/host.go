package glinstance

import (
	"fmt"
	"strings"
)

const defaultHostname = "gitlab.com"

var hostnameOverride string

// Default returns the host name of the default GitLab instance
func Default() string {
	return defaultHostname
}

// OverridableDefault is like Default, except it is overridable by the GITLAB_TOKEN environment variable
func OverridableDefault() string {
	if hostnameOverride != "" {
		return hostnameOverride
	}
	return Default()
}

// OverrideDefault overrides the value returned from OverridableDefault. This should only ever be
// called from the main runtime path, not tests.
func OverrideDefault(newhost string) {
	hostnameOverride = newhost
}

// IsSelfHosted reports whether a non-normalized host name looks like a Self-hosted GitLab instance
func IsSelfHosted(h string) bool {
	return NormalizeHostname(h) != Default()
}

// NormalizeHostname returns the canonical host name of a GitLab instance
// Taking cover in case GitLab allows subdomains on gitlab.com https://gitlab.com/gitlab-org/gitlab/-/issues/26703
func NormalizeHostname(h string) string {
	hostname := strings.ToLower(h)
	if strings.HasSuffix(hostname, "."+Default()) {
		return Default()
	}
	return hostname
}

// StripHostProtocol strips the url protocol and returns the hostname and the protocol
func StripHostProtocol(h string) (hostname, protocol string) {
	hostname = NormalizeHostname(h)
	if strings.HasPrefix(hostname, "http://") {
		protocol = "http"
	} else {
		protocol = "https"
	}
	hostname = strings.TrimPrefix(hostname, protocol)
	hostname = strings.Trim(hostname, "://")
	return
}

// APIEndpoint returns the API endpoint prefix for a GitLab instance :)
func APIEndpoint(hostname string) string {
	if IsSelfHosted(hostname) {
		return fmt.Sprintf("https://%s/api/v4/", hostname)
	}
	return "https://gitlab.com/api/v4/"
}
