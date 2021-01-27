.. _glab_mr_update:

glab mr update
--------------

Update merge requests

Synopsis
~~~~~~~~


Update merge requests

::

  glab mr update [<id> | <branch>] [flags]

Examples
~~~~~~~~

::

  $ glab mr update 23 --ready
  $ glab mr update 23 --draft
  $ glab mr update --draft  # Updates MR related to current branch
  

Options
~~~~~~~

::

  -a, --assignee strings       assign users via username, prefix with '!' or '-' to remove from existing assignees, '+' to add, otherwise replace existing assignees with given users
  -d, --description string     merge request description
      --draft                  Mark merge request as a draft
  -l, --label strings          add labels
      --lock-discussion        Lock discussion on merge request
  -m, --milestone string       title of the milestone to assign, pass "" or 0 to unassign
  -r, --ready                  Mark merge request as ready to be reviewed and merged
      --remove-source-branch   Remove Source Branch on merge
      --target-branch string   set target branch
  -t, --title string           Title of merge request
      --unassign               unassign all users
  -u, --unlabel strings        remove labels
      --unlock-discussion      Unlock discussion on merge request
      --wip                    Mark merge request as a work in progress. Alternative to --draft

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

