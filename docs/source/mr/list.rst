.. _glab_mr_list:

glab mr list
------------

List merge requests

Synopsis
~~~~~~~~


List merge requests

::

  glab mr list [flags]

Examples
~~~~~~~~

::

  $ glab mr list --all
  $ glab mr ls -a
  $ glab mr list --assignee=@me
  $ glab mr list --source-branch=new-feature
  $ glab mr list --target-branch=trunk
  $ glab mr list --search "this adds feature X"
  $ glab mr list --label needs-review
  $ glab mr list --not-label waiting-maintainer-feedback,subsystem-x
  $ glab mr list -M --per-page 10
  

Options
~~~~~~~

::

  -A, --all                    Get all merge requests
  -a, --assignee strings       Get only merge requests assigned to users
      --author string          Fitler merge request by Author <username>
  -c, --closed                 Get only closed merge requests
  -d, --draft                  Filter by draft merge requests
  -l, --label strings          Filter merge request by label <name>
  -M, --merged                 Get only merged merge requests
  -m, --milestone string       Filter merge request by milestone <id>
      --not-label strings      Filter merge requests by not having label <name>
  -p, --page int               Page number (default 1)
  -P, --per-page int           Number of items to list per page (default 30)
      --search string          Filter by <string> in title and description
  -s, --source-branch string   Filter by source branch <name>
  -t, --target-branch string   Filter by target branch <name>

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

