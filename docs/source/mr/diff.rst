.. _glab_mr_diff:

glab mr diff
------------

View changes in a merge request

Synopsis
~~~~~~~~


View changes in a merge request

::

  glab mr diff [<id> | <branch>] [flags]

Examples
~~~~~~~~

::

  $ glab mr diff 123
  $ glab mr diff branch
  $ glab mr diff  # get from current branch
  $ glab mr diff 123 --color=never
  

Options
~~~~~~~

::

      --color string   Use color in diff output: {always|never|auto} (default "auto")

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

