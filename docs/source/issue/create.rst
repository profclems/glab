.. _glab_issue_create:

glab issue create
-----------------

Create an issue

Synopsis
~~~~~~~~


Create an issue

::

  glab issue create [flags]

Options
~~~~~~~

::

  -a, --assignee string      Assign issue to people by their ID. Multiple values should be comma separated 
  -c, --confidential         Set an issue to be confidential. Default is false
  -d, --description string   Supply a description for issue
  -l, --label string         Add label by name. Multiple labels should be comma separated
      --linked-mr int        The IID of a merge request in which to resolve all issues (default -1)
  -m, --milestone int        The global ID of a milestone to assign issue (default -1)
      --no-editor            Don't open editor to enter description. If set to true, uses prompt. Default is false
  -t, --title string         Supply a title for issue
  -w, --weight int           The weight of the issue. Valid values are greater than or equal to 0. (default -1)

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

