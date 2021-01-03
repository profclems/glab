package glinstance

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestIsSelfHosted(t *testing.T) {
	type args struct {
		h string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "self_hosted subdomain",
			args: args{h: "gitlab.example.com"},
			want: true,
		},
		{
			name: "gitlab.com",
			args: args{h: "gitlab.com"},
			want: false,
		},
		{
			name: "is a gitlab.com subdomain",
			args: args{h: "example.gitlab.com"},
			want: false,
		},
		{
			name: "self hosted",
			args: args{h: "example.com"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSelfHosted(tt.args.h); got != tt.want {
				t.Errorf("IsSelfHosted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverridableDefault(t *testing.T) {
	oldOverride := hostnameOverride
	t.Cleanup(func() {
		hostnameOverride = oldOverride
	})

	host := OverridableDefault()
	if host != "gitlab.com" {
		t.Errorf("expected gitlab.com, got %q", host)
	}

	OverrideDefault("example.org")

	host = OverridableDefault()
	if host != "example.org" {
		t.Errorf("expected example.org, got %q", host)
	}
	host = Default()
	if host != "gitlab.com" {
		t.Errorf("expected gitlab.com, got %q", host)
	}
}

func TestNormalizeHostname(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{
			host: "GitLab.com",
			want: "gitlab.com",
		},
		{
			host: "subdomain.gitlab.com",
			want: "gitlab.com",
		},
		{
			host: "ssh.gitlab.com",
			want: "gitlab.com",
		},
		{
			host: "upload.gitlab.com",
			want: "gitlab.com",
		},
		{
			host: "EXAMPLE.COM",
			want: "example.com",
		},
		{
			host: "gitlab.my.org",
			want: "gitlab.my.org",
		},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			if got := NormalizeHostname(tt.host); got != tt.want {
				t.Errorf("NormalizeHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIEndpoint(t *testing.T) {
	tests := []struct {
		host     string
		protocol string
		want     string
	}{
		{
			host: "gitlab.com",
			want: "https://gitlab.com/api/v4/",
		},
		{
			host: "ghe.io",
			want: "https://ghe.io/api/v4/",
		},
		{
			host:     "salsa.debian.com",
			protocol: "http",
			want:     "http://salsa.debian.com/api/v4/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			if got := APIEndpoint(tt.host, tt.protocol); got != tt.want {
				t.Errorf("APIEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "default_hostname",
			want: "gitlab.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); got != tt.want {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripHostProtocol(t *testing.T) {
	tests := []struct {
		name         string
		hostname     string
		wantHostname string
		wantProtocol string
	}{
		{
			name:         "url with https protocol",
			hostname:     "https://gitlab.com",
			wantHostname: "gitlab.com",
			wantProtocol: "https",
		},
		{
			name:         "https url with ending slash",
			hostname:     "https://gitlab.com/",
			wantHostname: "gitlab.com",
			wantProtocol: "https",
		},
		{
			name:         "url with http protocol",
			hostname:     "http://gitlab.com/",
			wantHostname: "gitlab.com",
			wantProtocol: "http",
		},
		{
			name:         "http url with ending slash",
			hostname:     "http://gitlab.com/",
			wantHostname: "gitlab.com",
			wantProtocol: "http",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostname, gotProtocol := StripHostProtocol(tt.hostname)
			if gotHostname != tt.wantHostname {
				t.Errorf("StripHostProtocol() gotHostname = %v, want %v", gotHostname, tt.wantHostname)
			}
			if gotProtocol != tt.wantProtocol {
				t.Errorf("StripHostProtocol() gotProtocol = %v, want %v", gotProtocol, tt.wantProtocol)
			}
		})
	}
}

func Test(t *testing.T) {
	testCases := []struct {
		name     string
		hostname interface{}
		expected string
	}{
		{
			name:     "valid",
			hostname: "localhost",
		},
		{
			name:     "invalid/not-string",
			hostname: 1,
			expected: "hostname is not a string",
		},
		{
			name:     "invalid/empty-string",
			hostname: "",
			expected: "a value is required",
		},
		{
			name:     "invalid/has-foward-slash",
			hostname: "local/host",
			expected: "invalid hostname",
		},
		{
			name:     "invalid/has-colon",
			hostname: "local:host",
			expected: "invalid hostname",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			err := HostnameValidator(tC.hostname)
			if tC.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tC.expected)
			}
		})
	}
}

func Test_GraphQLEndpoint(t *testing.T) {
	testCases := []struct {
		name     string
		protocol string
		hostname string
		output   string
	}{
		{
			name:     "OfficialInstance/https",
			protocol: "https",
			hostname: "gitlab.com",
			output:   "https://gitlab.com/api/graphql/",
		},
		{
			name:     "OfficialInstance/any-protocol-is-https",
			protocol: "NoExistProtocol",
			hostname: "gitlab.com",
			output:   "https://gitlab.com/api/graphql/",
		},
		{
			name:     "OfficialInstance/no-protocol-default-to-https",
			protocol: "",
			hostname: "gitlab.com",
			output:   "https://gitlab.com/api/graphql/",
		},
		{
			name:     "SelfHosted/https",
			protocol: "https",
			hostname: "gitlab.alpinelinux.org",
			output:   "https://gitlab.alpinelinux.org/api/graphql/",
		},
		{
			name:     "SelfHost/http",
			protocol: "http",
			hostname: "gitlab.alpinelinux.org",
			output:   "http://gitlab.alpinelinux.org/api/graphql/",
		},
		{
			name:     "SelfHosted/no-protocol-default-to-https",
			protocol: "",
			hostname: "gitlab.alpinelinux.org",
			output:   "https://gitlab.alpinelinux.org/api/graphql/",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := GraphQLEndpoint(tC.hostname, tC.protocol)
			assert.Equal(t, tC.output, got)
		})
	}
}
