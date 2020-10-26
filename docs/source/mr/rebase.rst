.. _glab_mr_rebase:

glab mr rebase
--------------

Automatically rebase the source_branch of the merge request against its target_branch.

Synopsis
~~~~~~~~


If you don’t have permissions to push to the merge request’s source branch - you’ll get a 403 Forbidden response.

::

  glab mr rebase <id> [flags]

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

