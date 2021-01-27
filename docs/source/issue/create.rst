.. _glab_issue_create:

glab issue create
-----------------

Create an issue

Synopsis
~~~~~~~~


Create an issue

::

  glab issue create [flags]

Examples
~~~~~~~~

::

  $ glab issue create
  $ glab issue new
  $ glab issue create -m release-2.0.0 -t "we need this feature" --label important
  $ glab issue new -t "Fix CVE-YYYY-XXXX" -l security --linked-mr 123
  $ glab issue create -m release-1.0.1 -t "security fix" --label security --web
  

Options
~~~~~~~

::

  -a, --assignee usernames   Assign issue to people by their usernames
  -c, --confidential         Set an issue to be confidential. Default is false
  -d, --description string   Supply a description for issue
  -l, --label strings        Add label by name. Multiple labels should be comma separated
      --linked-mr int        The IID of a merge request in which to resolve all issues
  -m, --milestone string     The global ID or title of a milestone to assign
      --no-editor            Don't open editor to enter description. If set to true, uses prompt. Default is false
  -t, --title string         Supply a title for issue
      --web                  continue issue creation with web interface
  -w, --weight int           The weight of the issue. Valid values are greater than or equal to 0.
  -y, --yes                  Don't prompt for confirmation to submit the issue

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

