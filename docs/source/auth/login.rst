.. _glab_auth_login:

glab auth login
---------------

Authenticate with a GitLab instance

Synopsis
~~~~~~~~


Authenticate with a GitLab instance.
You can pass in a token on standard input by using `--stdin`.
The minimum required scopes for the token are: "api", "write_repository".


::

  glab auth login [flags]

Examples
~~~~~~~~

::

  # start interactive setup
  $ glab auth login
  # authenticate against gitlab.com by reading the token from a file
  $ glab auth login --stdin < myaccesstoken.txt
  # authenticate with a self-hosted GitLab instance
  $ glab auth login --hostname salsa.debian.org
  

Options
~~~~~~~

::

  -h, --hostname string   The hostname of the GitLab instance to authenticate with
      --stdin             Read token from standard input
  -t, --token string      Your GitLab access token

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

