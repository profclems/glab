.. _glab_mr_delete:

glab mr delete
--------------

Delete merge requests

Synopsis
~~~~~~~~


Delete merge requests

::

  glab mr delete [<id> | <branch>] [flags]

Examples
~~~~~~~~

::

  $ glab mr delete 123
  $ glab mr delete 123 branch-name 789  # close multiple branches
  $ glab mr delete 1,2,branch-related-to-mr-3,4,5  # close MRs !1,!2,!3,!4,!5
  $ glab mr del 123
  $ glab mr delete branch
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

