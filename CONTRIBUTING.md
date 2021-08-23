# Contributing

[legal]: https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license
[license]: LICENSE

Hi! Thanks for your interest in contributing to this project!

To encourage active collaboration, pull requests are strongly encouraged, not just bug reports. "Bug reports" may also be sent in the form of a pull request containing a failing test. I'd also love to hear about ideas for new features as issues.

Please do:

* Check existing issues to verify that the bug or feature request has not already been submitted.
* Open an issue if things aren't working as expected.
* Open an issue to propose a significant change.
* open an issue to propose a feature
* Open a pull request to fix a bug.
* Open a pull request to fix documentation about a command.
* Open a pull request for an issue with the help-wanted label and leave a comment claiming it.

Please avoid:

* Opening pull requests for issues marked `needs-design`, `needs-investigation`, `needs-user-input`, or `blocked`.
* Opening pull requests for documentation for a new command specifically. Manual pages are auto-generated from source after every release

## Building the project

Prerequisites:
- Go 1.13+

Build with: `make` or `go build -o bin/glab ./cmd/glab/main.go`

Run the new binary as: `./bin/glab`

Run tests with: `make test` or `go test ./...`

> WARNING: Do not run `make test` outside of an isolated environment, it will overwrite your global config.

## Submitting a pull request

1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and ensure tests pass
1. Submit a pull request

## Commits Message

### TL;DR: Your commit message should be semantic

A commit message consists of a header, a body and a footer, separated by a blank line.

Any line of the commit message cannot be longer than 100 characters! This allows the message to be easier to read on GitHub as well as in various git tools.

```sh
<type>[optional scope]: <description>
<BLANK LINE>
[optional body]
<BLANK LINE>
<footer>
```

### Message Header
Ideally, the commit message heading which contains the description, should not be more than 50 characters

The message header is a single line that contains a succinct description of the change containing a type, an optional scope, and a subject.

#### `<type>`

This describes the kind of change that this commit is providing

- feat (feature)
- fix (bug fix)
- docs (documentation)
- style (formatting, missing semicolons, …)
- refactor(restructuring codebase)
- test (when adding missing tests)
- chore (maintain)

#### `<scope>`

Scope can be anything specifying the place of the commit change. For example events, kafka, userModel, authorization, authentication, loginPage, etc

#### `<description>`

This is a very short description of the change

* `use imperative, present tense: “change” not “changed” nor “changes”`
* `don't capitalize the first letter`
* `no dot (.) at the end`

## Message Body

- just as in subject use imperative, present tense: “change” not “changed” nor “changes”
- includes motivation for the change and contrasts with previous behavior

<http://365git.tumblr.com/post/3308646748/writing-git-commit-messages>

<http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html>

## Message Footer

Finished, fixed or delivered stories should be listed on a separate line in the footer prefixed with "Finishes", "Fixes" , or "Delivers" keyword like this:

`[(Finishes|Fixes|Delivers) #ISSUE_ID]`

## Message Example

```sh
feat(kafka): implement exactly once delivery

- ensure every event published to kafka is delivered exactly once
- implement error handling for failed delivery

Delivers #065
```

```sh
fix(login): allow provided user preferences to override default preferences

- This allows the preferences associated with a user account to override and customize the default app preference like theme, timezone e.t.c

Fixes #025
```

Contributions to this project are made available to public under the [project's open source license][license].
Please note that this project adheres to a [Contributor Code of Conduct](https://github.com/profclems/glab/tree/trunk/CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

Manual pages are auto-generated from source on every release. You do not need to submit pull requests for documentation specifically; manual pages for commands will automatically get updated after your pull requests gets accepted.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
