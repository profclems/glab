.. _glab_mr_update:

glab mr update
--------------

Update merge requests

Synopsis
~~~~~~~~


Update merge requests

::

  glab mr update <id> [flags]

Examples
~~~~~~~~

::

  $ glab mr update 23 --ready
  $ glab mr update 23 --draft
  $ glab mr update --draft  # Updates MR related to current branch
  

Options
~~~~~~~

::

  -a, --assignee string        merge request assignee
  -d, --description string     merge request description
      --draft                  Mark merge request as a draft
      --lock-discussion        Lock discussion on merge request
  -r, --ready                  Mark merge request as ready to be reviewed and merged
      --remove-source-branch   Remove Source Branch on merge
  -t, --title string           Title of merge request
      --wip                    Mark merge request as a work in progress. Alternative to --draft

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

