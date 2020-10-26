.. _glab_mr_for:

glab mr for
-----------

Create new merge request for an issue

Synopsis
~~~~~~~~


Create new merge request for an issue

::

  glab mr for [flags]

Examples
~~~~~~~~

::

  $ glab mr for 34   # Create mr for issue 34
  $ glab mr for 34 --wip   # Create mr and mark as work in progress
  $ glab mr new-for 34
  $ glab mr create-for 34
  

Options
~~~~~~~

::

      --allow-collaboration    Allow commits from other members
  -a, --assignee string        Assign merge request to people by their IDs. Multiple values should be comma separated 
      --draft                  Mark merge request as a draft. Default is true (default true)
  -l, --label string           Add label by name. Multiple labels should be comma separated
  -m, --milestone int          add milestone by <id> for merge request (default -1)
      --remove-source-branch   Remove Source Branch on merge
  -b, --target-branch string   The target or base branch into which you want your code merged
      --wip                    Mark merge request as a work in progress. Overrides --draft

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

