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

  -a, --assignee strings     assign users via username, prefix with '!' or '-' to remove from existing assignees, '+' to add, otherwise replace existing assignees with given users
  -c, --confidential         Make issue confidential
  -d, --description string   Issue description
  -l, --label strings        add labels
      --lock-discussion      Lock discussion on issue
  -m, --milestone string     title of the milestone to assign, pass "" or 0 to unassign
  -p, --public               Make issue public
  -t, --title string         Title of issue
      --unassign             unassign all users
  -u, --unlabel strings      remove labels
      --unlock-discussion    Unlock discussion on issue

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

