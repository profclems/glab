.. _glab_mr_merge:

glab mr merge
-------------

Merge/Accept merge requests

Synopsis
~~~~~~~~


Merge/Accept merge requests

::

  glab mr merge {<id> | <branch>} [flags]

Examples
~~~~~~~~

::

  $ glab mr merge 235
  $ glab mr accept 235
  $ glab mr merge    # Finds open merge request from current branch
  

Options
~~~~~~~

::

  -m, --message string           Custom merge commit message
  -r, --rebase                   Rebase the commits onto the base branch
  -d, --remove-source-branch     Remove source branch on merge
      --sha string               Merge Commit sha
  -s, --squash                   Squash commits on merge
      --squash-message string    Custom Squash commit message
      --when-pipeline-succeeds   Merge only when pipeline succeeds (default true)
  -y, --yes                      Skip submission confirmation prompt

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

