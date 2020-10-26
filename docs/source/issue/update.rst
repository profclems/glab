.. _glab_issue_update:

glab issue update
-----------------

Update issue

Synopsis
~~~~~~~~


Update issue

::

  glab issue update <id> [flags]

Examples
~~~~~~~~

::

  $ glab issue update 42 --label ui,ux
  $ glab issue update 42 --unlabel working
  

Options
~~~~~~~

::

  -d, --description string    Issue description
  -l, --label stringArray     add labels
      --lock-discussion       Lock discussion on issue
  -t, --title string          Title of issue
  -u, --unlabel stringArray   remove labels

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

