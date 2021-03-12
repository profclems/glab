.. _glab_mr_reopen:

glab mr reopen
--------------

Reopen merge requests

Synopsis
~~~~~~~~


Reopen merge requests

::

  glab mr reopen [<id>... | <branch>...] [flags]

Examples
~~~~~~~~

::

  $ glab mr reopen 123
  $ glab mr reopen 123 456 789
  $ glab mr reopen branch-1 branch-2
  $ glab mr reopen  # use checked out branch
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

