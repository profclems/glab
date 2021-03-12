.. _glab_mr_revoke:

glab mr revoke
--------------

Revoke approval on a merge request <id>

Synopsis
~~~~~~~~


Revoke approval on a merge request <id>

::

  glab mr revoke [<id> | <branch>] [flags]

Examples
~~~~~~~~

::

  $ glab mr revoke 123
  $ glab mr unapprove 123
  $ glab mr revoke branch
  $ glab mr revoke  # use checked out branch
  $ glab mr revoke 123 branch 456
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

