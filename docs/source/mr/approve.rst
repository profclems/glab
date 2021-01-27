.. _glab_mr_approve:

glab mr approve
---------------

Approve merge requests

Synopsis
~~~~~~~~


Approve merge requests

::

  glab mr approve {<id> | <branch>} [flags]

Examples
~~~~~~~~

::

  $ glab mr approve 235
  $ glab mr approve    # Finds open merge request from current branch
  

Options
~~~~~~~

::

  -s, --sha string   SHA which must match the SHA of the HEAD commit of the merge request

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

