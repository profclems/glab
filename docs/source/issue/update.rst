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

  -c, --confidential          Make issue confidential
  -d, --description string    Issue description
  -l, --label stringArray     add labels
      --lock-discussion       Lock discussion on issue
  -p, --public                Make issue public
  -t, --title string          Title of issue
  -u, --unlabel stringArray   remove labels
      --unlock-discussion     Unlock discussion on issue

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

