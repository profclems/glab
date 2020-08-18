## Contributing

[legal]: https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license
[license]: ../LICENSE

Hi! Thanks for your interest in contributing to this project!

To encourage active collaboration, pull requests are strongly encouraged, not just bug reports. "Bug reports" may also be sent in the form of a pull request containing a failing test.

Please do:

* open an issue if things aren't working as expected
* open an issue to propose a significant change
* open an issue to propose a feature
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

## Commits Message
A commit message consists of a header, a body and a footer, separated by a blank line.

Any line of the commit message cannot be longer than 100 characters! This allows the message to be easier to read on GitHub as well as in various git tools.

```
<type>[optional scope]: <description> 
<BLANK LINE>
[optional body]
<BLANK LINE>
<footer>
```

### Message Header
Ideally, the commit message heading which contains the description, should not be more than 50 characters

The message header is a single line that contains a succinct description of the change containing a type, an optional scope, and a subject.

#####`<type>` This describes the kind of change that this commit is providing.

- feat (feature)
- fix (bug fix)
- docs (documentation)
- style (formatting, missing semicolons, …)
- refactor(restructuring codebase)
- test (when adding missing tests)
- chore (maintain)

#####`<scope>` Scope can be anything specifying the place of the commit change. For example events, kafka, userModel, authorization, authentication, loginPage, etc...

#####`<subject>` This is a very short description of the change.

* `use imperative, present tense: “change” not “changed” nor “changes”`
* `don't capitalize the first letter`
* `no dot (.) at the end`

## Message Body

- just as in subject use imperative, present tense: “change” not “changed” nor “changes”
- includes motivation for the change and contrasts with previous behavior

http://365git.tumblr.com/post/3308646748/writing-git-commit-messages

http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html

## Message Footer

Finished, fixed or delivered stories should be listed on a separate line in the footer prefixed with "Finishes", "Fixes" , or "Delivers" keyword like this:

`[(Finishes|Fixes|Delivers) #ISSUE_ID]`

## Message Example

```
feat(kafka): implement exactly once delivery

- ensure every event published to kafka is delivered exactly once
- implement error handling for failed delivery

Delivers #065
```
```
fix(login): allow provided user preferences to override default preferences

- This allows the preferences associated with a user account to override and customize the default app preference like theme, timezone e.t.c

Fixes #025
```


Contributions to this project are made available to public under the [project's open source license][license].

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
