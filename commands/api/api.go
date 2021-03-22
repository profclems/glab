package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/spf13/cobra"
	jsonPretty "github.com/tidwall/pretty"
	"github.com/xanzy/go-gitlab"
)

type ApiOptions struct {
	IO *iostreams.IOStreams

	HttpClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Branch     func() (string, error)
	Config     config.Config

	Hostname            string
	RequestMethod       string
	RequestMethodPassed bool
	RequestPath         string
	RequestInputFile    string
	MagicFields         []string
	RawFields           []string
	RequestHeaders      []string
	ShowResponseHeaders bool
	Paginate            bool
	Silent              bool
}

func NewCmdApi(f *cmdutils.Factory, runF func(*ApiOptions) error) *cobra.Command {
	opts := ApiOptions{
		IO:         f.IO,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
		Branch:     f.Branch,
	}

	cmd := &cobra.Command{
		Use:   "api <endpoint>",
		Short: "Make an authenticated request to GitLab API",
		Long: `Makes an authenticated HTTP request to the GitLab API and prints the response.
The endpoint argument should either be a path of a GitLab API v4 endpoint, or 
"graphql" to access the GitLab's GraphQL API.

GitLab REST API Docs: https://docs.gitlab.com/ce/api/README.html
GitLab GraphQL Docs: https://docs.gitlab.com/ee/api/graphql/

If the current directory is a git directory, the GitLab authenticated host in the current git 
directory will be used else gitlab.com will be used.
Override the GitLab hostname with '--hostname'.

Placeholder values ":fullpath" or ":id"", ":user" or ":username", ":group", ":namespace", 
":repo", and ":branch" in the endpoint argument will get replaced with values from the 
repository of the current directory.

The default HTTP request method is "GET" normally and "POST" if any parameters
were added. Override the method with '--method'.

Pass one or more '--raw-field' values in "key=value" format to add
JSON-encoded string parameters to the POST body.

The '--field' flag behaves like '--raw-field' with magic type conversion based
on the format of the value:
- literal values "true", "false", "null", and integer numbers get converted to
  appropriate JSON types;
- placeholder values ":namespace", ":repo", and ":branch" get populated with values
  from the repository of the current directory;
- if the value starts with "@", the rest of the value is interpreted as a
  filename to read the value from. Pass "-" to read from standard input.

For GraphQL requests, all fields other than "query" and "operationName" are
interpreted as GraphQL variables.

Raw request body may be passed from the outside via a file specified by '--input'.
Pass "-" to read from standard input. In this mode, parameters specified via
'--field' flags are serialized into URL query parameters.

In '--paginate' mode, all pages of results will sequentially be requested until
there are no more pages of results. For GraphQL requests, this requires that the
original query accepts an '$endCursor: String' variable and that it fetches the
'pageInfo{ hasNextPage, endCursor }' set of fields from a collection.`,
		Example: heredoc.Doc(`
			$ glab api projects/:fullpath/releases

			$ glab api projects/gitlab-com%2Fwww-gitlab-com/issues

			$ glab api issues --paginate

			$ glab api graphql -f query='
			  query {
			    project(fullPath: "gitlab-org/gitlab-docs") {
			      name
			      forksCount
			      statistics {
			        wikiSize
			      }
			      issuesEnabled
			      boards {
			        nodes {
			          id
			          name
			        }
			      }
			    }
			  }
			'

			$ glab api graphql --paginate -f query='
			  query($endCursor: String) {
			    project(fullPath: "gitlab-org/graphql-sandbox") {
			      name
			      issues(first: 2, after: $endCursor) {
			        edges {
			          node {
			            title
			          }
			        }
			        pageInfo {
			          endCursor
			          hasNextPage
			        }
			      }
			    }
			  }'
		`),
		Annotations: map[string]string{
			"help:environment": heredoc.Doc(`
				GITLAB_TOKEN, OAUTH_TOKEN (in order of precedence): an authentication token for API requests.
				GITLAB_HOST, GITLAB_URI, GITLAB_URL: specify a GitLab host to make request to.
			`),
		},
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			opts.RequestPath = args[0]
			opts.RequestMethodPassed = c.Flags().Changed("method")
			opts.Config, _ = f.Config()

			if c.Flags().Changed("hostname") {
				if err := glinstance.HostnameValidator(opts.Hostname); err != nil {
					return &cmdutils.FlagError{Err: fmt.Errorf("error parsing --hostname: %w", err)}
				}
			}

			if opts.Paginate && !strings.EqualFold(opts.RequestMethod, "GET") && opts.RequestPath != "graphql" {
				return &cmdutils.FlagError{Err: errors.New(`the '--paginate' option is not supported for non-GET requests`)}
			}
			if opts.Paginate && opts.RequestInputFile != "" {
				return &cmdutils.FlagError{Err: errors.New(`the '--paginate' option is not supported with '--input'`)}
			}

			if runF != nil {
				return runF(&opts)
			}
			return apiRun(&opts)
		},
	}

	cmd.Flags().StringVar(&opts.Hostname, "hostname", "", "The GitLab hostname for the request (default is \"gitlab.com\" or authenticated host in current git directory)")
	cmd.Flags().StringVarP(&opts.RequestMethod, "method", "X", "GET", "The HTTP method for the request")
	cmd.Flags().StringArrayVarP(&opts.MagicFields, "field", "F", nil, "Add a parameter of inferred type")
	cmd.Flags().StringArrayVarP(&opts.RawFields, "raw-field", "f", nil, "Add a string parameter")
	cmd.Flags().StringArrayVarP(&opts.RequestHeaders, "header", "H", nil, "Add an additional HTTP request header")
	cmd.Flags().BoolVarP(&opts.ShowResponseHeaders, "include", "i", false, "Include HTTP response headers in the output")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Make additional HTTP requests to fetch all pages of results")
	cmd.Flags().StringVar(&opts.RequestInputFile, "input", "", "The file to use as body for the HTTP request")
	cmd.Flags().BoolVar(&opts.Silent, "silent", false, "Do not print the response body")
	return cmd
}

func apiRun(opts *ApiOptions) error {
	params, err := parseFields(opts)
	if err != nil {
		return err
	}
	isGraphQL := opts.RequestPath == "graphql"
	requestPath, err := fillPlaceholders(opts.RequestPath, opts)
	if err != nil {
		return fmt.Errorf("unable to expand placeholder in path: %w", err)
	}
	method := opts.RequestMethod
	requestHeaders := opts.RequestHeaders
	var requestBody interface{} = params

	if !opts.RequestMethodPassed && (len(params) > 0 || opts.RequestInputFile != "") {
		method = "POST"
	}

	if opts.Paginate && !isGraphQL {
		requestPath = addPerPage(requestPath, 100, params)
	}

	if opts.RequestInputFile != "" {
		file, size, err := openUserFile(opts.RequestInputFile, opts.IO.In)
		if err != nil {
			return err
		}
		defer file.Close()
		requestPath, err = parseQuery(requestPath, params)
		if err != nil {
			return err
		}
		requestBody = file
		if size >= 0 {
			requestHeaders = append([]string{fmt.Sprintf("Content-Length: %d", size)}, requestHeaders...)
		}
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return err
	}

	headersOutputStream := opts.IO.StdOut
	if opts.Silent {
		opts.IO.StdOut = ioutil.Discard
	} else {
		err := opts.IO.StartPager()
		if err != nil {
			return err
		}
		defer opts.IO.StopPager()
	}

	host := httpClient.BaseURL().Host
	if opts.Hostname != "" {
		host = opts.Hostname
	}

	hasNextPage := true
	for hasNextPage {
		resp, err := httpRequest(api.GetClient(), opts.Config, host, method, requestPath, requestBody, requestHeaders)
		if err != nil {
			return err
		}

		endCursor, err := processResponse(resp, opts, headersOutputStream)
		if err != nil {
			return err
		}

		if !opts.Paginate {
			break
		}

		if isGraphQL {
			hasNextPage = endCursor != ""
			if hasNextPage {
				params["endCursor"] = endCursor
			}
		} else {
			requestPath, hasNextPage = findNextPage(resp)
		}

		if hasNextPage && opts.ShowResponseHeaders {
			fmt.Fprint(opts.IO.StdOut, "\n")
		}
	}

	return nil
}

func processResponse(resp *http.Response, opts *ApiOptions, headersOutputStream io.Writer) (endCursor string, err error) {
	if opts.ShowResponseHeaders {
		fmt.Fprintln(headersOutputStream, resp.Proto, resp.Status)
		printHeaders(headersOutputStream, resp.Header, opts.IO.ColorEnabled())
		fmt.Fprint(headersOutputStream, "\r\n")
	}

	if resp.StatusCode == 204 {
		return
	}
	var responseBody io.Reader = resp.Body

	isJSON, _ := regexp.MatchString(`[/+]json(;|$)`, resp.Header.Get("Content-Type"))

	var serverError string
	if isJSON && (opts.RequestPath == "graphql" || resp.StatusCode >= 400) {
		responseBody, serverError, err = parseErrorResponse(responseBody, resp.StatusCode)
		if err != nil {
			return
		}
	}

	var bodyCopy *bytes.Buffer
	isGraphQLPaginate := isJSON && resp.StatusCode == 200 && opts.Paginate && opts.RequestPath == "graphql"
	if isGraphQLPaginate {
		bodyCopy = &bytes.Buffer{}
		responseBody = io.TeeReader(responseBody, bodyCopy)
	}

	if isJSON && opts.IO.ColorEnabled() {
		out := &bytes.Buffer{}
		_, err = io.Copy(out, responseBody)
		if err == nil {
			result := jsonPretty.Color(jsonPretty.Pretty(out.Bytes()), nil)
			_, err = fmt.Fprintln(opts.IO.StdOut, string(result))
		}
	} else {
		_, err = io.Copy(opts.IO.StdOut, responseBody)
	}
	if err != nil {
		return
	}

	if serverError != "" {
		fmt.Fprintf(opts.IO.StdErr, "glab: %s\n", serverError)
		err = cmdutils.SilentError
		return
	} else if resp.StatusCode > 299 {
		fmt.Fprintf(opts.IO.StdErr, "glab: HTTP %d\n", resp.StatusCode)
		err = cmdutils.SilentError
		return
	}

	if isGraphQLPaginate {
		endCursor = findEndCursor(bodyCopy)
	}

	return
}

var placeholderRE = regexp.MustCompile(`:(group/:namespace/:repo|namespace/:repo|fullpath|id|user|username|group|namespace|repo|branch)\b`)

// fillPlaceholders populates `:namespace` and `:repo` placeholders with values from the current repository
func fillPlaceholders(value string, opts *ApiOptions) (string, error) {
	if !placeholderRE.MatchString(value) {
		return value, nil
	}

	baseRepo, err := opts.BaseRepo()
	if err != nil {
		return value, err
	}

	filled := placeholderRE.ReplaceAllStringFunc(value, func(m string) string {
		switch m {
		case ":id":
			h, _ := opts.HttpClient()
			baseProject, e := api.GetProject(h, baseRepo.FullName())
			if e == nil && baseProject != nil {
				return string(rune(baseProject.ID))
			}
			err = e
			return ""
		case ":group/:namespace/:repo", ":fullpath":
			return url.PathEscape(baseRepo.FullName())
		case ":namespace/:repo":
			return url.PathEscape(baseRepo.RepoNamespace() + "/" + baseRepo.RepoName())
		case ":group":
			return baseRepo.RepoGroup()
		case ":user", ":username":
			h, _ := opts.HttpClient()
			u, e := api.CurrentUser(h)
			if e == nil && u != nil {
				return u.Username
			}
			err = e
			return m
		case ":namespace":
			return baseRepo.RepoNamespace()
		case ":repo":
			return baseRepo.RepoName()
		case ":branch":
			branch, e := opts.Branch()
			if e != nil {
				err = e
			}
			return branch
		default:
			err = fmt.Errorf("invalid placeholder: %q", m)
			return ""
		}
	})

	if err != nil {
		return value, err
	}

	return filled, nil
}

func printHeaders(w io.Writer, headers http.Header, colorize bool) {
	var names []string
	for name := range headers {
		if name == "Status" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	var headerColor, headerColorReset string
	if colorize {
		headerColor = "\x1b[1;34m" // bright blue
		headerColorReset = "\x1b[m"
	}
	for _, name := range names {
		fmt.Fprintf(w, "%s%s%s: %s\r\n", headerColor, name, headerColorReset, strings.Join(headers[name], ", "))
	}
}

func parseFields(opts *ApiOptions) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	for _, f := range opts.RawFields {
		key, value, err := parseField(f)
		if err != nil {
			return params, err
		}
		params[key] = value
	}
	for _, f := range opts.MagicFields {
		key, strValue, err := parseField(f)
		if err != nil {
			return params, err
		}
		value, err := magicFieldValue(strValue, opts)
		if err != nil {
			return params, fmt.Errorf("error parsing %q value: %w", key, err)
		}
		params[key] = value
	}
	return params, nil
}

func parseField(f string) (string, string, error) {
	idx := strings.IndexRune(f, '=')
	if idx == -1 {
		return f, "", fmt.Errorf("field %q requires a value separated by an '=' sign", f)
	}
	return f[0:idx], f[idx+1:], nil
}

func magicFieldValue(v string, opts *ApiOptions) (interface{}, error) {
	if strings.HasPrefix(v, "@") {
		return readUserFile(v[1:], opts.IO.In)
	}

	if n, err := strconv.Atoi(v); err == nil {
		return n, nil
	}

	switch v {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	default:
		return fillPlaceholders(v, opts)
	}
}

func readUserFile(fn string, stdin io.ReadCloser) ([]byte, error) {
	var r io.ReadCloser
	if fn == "-" {
		r = stdin
	} else {
		var err error
		r, err = os.Open(fn)
		if err != nil {
			return nil, err
		}
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func openUserFile(fn string, stdin io.ReadCloser) (io.ReadCloser, int64, error) {
	if fn == "-" {
		return stdin, -1, nil
	}

	r, err := os.Open(fn)
	if err != nil {
		return r, -1, err
	}

	s, err := os.Stat(fn)
	if err != nil {
		return r, -1, err
	}

	return r, s.Size(), nil
}

func parseErrorResponse(r io.Reader, statusCode int) (io.Reader, string, error) {
	bodyCopy := &bytes.Buffer{}
	b, err := ioutil.ReadAll(io.TeeReader(r, bodyCopy))
	if err != nil {
		return r, "", err
	}

	var parsedBody struct {
		Message string
		Errors  []json.RawMessage
	}
	err = json.Unmarshal(b, &parsedBody)
	if err != nil {
		return r, "", err
	}
	if parsedBody.Message != "" {
		return bodyCopy, fmt.Sprintf("%s (HTTP %d)", parsedBody.Message, statusCode), nil
	}

	type errorMessage struct {
		Message string
	}
	var respErrors []string
	for _, rawErr := range parsedBody.Errors {
		if len(rawErr) == 0 {
			continue
		}
		if rawErr[0] == '{' {
			var objectError errorMessage
			err := json.Unmarshal(rawErr, &objectError)
			if err != nil {
				return r, "", err
			}
			respErrors = append(respErrors, objectError.Message)
		} else if rawErr[0] == '"' {
			var stringError string
			err := json.Unmarshal(rawErr, &stringError)
			if err != nil {
				return r, "", err
			}
			respErrors = append(respErrors, stringError)
		}
	}

	if len(respErrors) > 0 {
		return bodyCopy, strings.Join(respErrors, "\n"), nil
	}

	return bodyCopy, "", nil
}
