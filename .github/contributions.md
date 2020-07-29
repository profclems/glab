## Contributing

[legal]: https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license
[license]: ../LICENSE

Hi! Thanks for your interest in contributing to this project!

To encourage active collaboration, pull requests are strongly encounraged, not just bug reports. "Bug reports" may also be sent in the form of a pull request containing a failing test.

Please do:

* open an issue if things aren't working as expected
* open an issue to propose a significant change
* open an issue to propose a fearure
* open a pull request to fix a bug
* open a pull request to fix documentation about a command
* open a pull request if your issue is marked as relevant by a community member after having discussed the issue

## Building the project

Prerequisites:
- Go 1.13

Build with: `make build` or `go build -o bin/glab ./cmd/main.go`

Run the new binary as: `./bin/glab`

## Submitting a pull request

1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change
1. Submit a pull request

## Branch Naming
Branches created should be named using the following format:

`
{story type}-{2-5 word summary}
`
`Issue or story type prefixes:` Indicates the context of the branch and should be one of:
- ft == Feature
- ch == Chore
- bg == Bug
- rf == Refractor

`Story Summary` -  Short 2-5 words summary about what the branch contains

### Example
`ft-gitlab-auth`

`bg-gitlab-auth-fails`


Contributions to this project are made available to public under the [project's open source license][license].

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
