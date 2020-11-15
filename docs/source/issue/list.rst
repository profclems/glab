.. _glab_issue_list:

glab issue list
---------------

List project issues

Synopsis
~~~~~~~~


List project issues

::

  glab issue list [flags]

Options
~~~~~~~

::

  -a, --all                Get all issues
      --assignee string    Filter issue by assignee <username>
  -c, --closed             Get only closed issues
      --confidential       Filter by confidential issues
  -l, --label string       Filter issue by label <name>
      --milestone string   Filter issue by milestone <id>
      --mine               Filter only issues issues assigned to me
  -o, --opened             Get only opened issues
  -p, --page int           Page number (default 1)
  -P, --per-page int       Number of items to list per page. (default 30) (default 30)

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

