.. _glab_auth_status:

glab auth status
----------------

View authentication status

Synopsis
~~~~~~~~


Verifies and displays information about your authentication state.

This command tests the authentication states of all known GitLab instances in the config file and reports issues if any


::

  glab auth status [flags]

Options
~~~~~~~

::

  -h, --hostname string   Check a specific instance's authentication status
  -t, --show-token        Display the auth token

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

