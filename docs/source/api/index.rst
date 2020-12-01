.. _glab_api:

glab api
--------

Make an authenticated request to GitLab API

Synopsis
~~~~~~~~


Makes an authenticated HTTP request to the GitLab API and prints the response.
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
'pageInfo{ hasNextPage, endCursor }' set of fields from a collection.

::

  glab api <endpoint> [flags]

Examples
~~~~~~~~

::

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
  

Options
~~~~~~~

::

  -F, --field stringArray       Add a parameter of inferred type
  -H, --header stringArray      Add an additional HTTP request header
      --hostname string         The GitLab hostname for the request (default is "gitlab.com" or authenticated host in current git directory)
  -i, --include                 Include HTTP response headers in the output
      --input string            The file to use as body for the HTTP request
  -X, --method string           The HTTP method for the request (default "GET")
      --paginate                Make additional HTTP requests to fetch all pages of results
  -f, --raw-field stringArray   Add a string parameter
      --silent                  Do not print the response body

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

