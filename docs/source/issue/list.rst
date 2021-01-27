.. _glab_issue_list:

glab issue list
---------------

List project issues

Synopsis
~~~~~~~~


List project issues

::

  glab issue list [flags]

Examples
~~~~~~~~

::

  $ glab issue list --all
  $ glab issue ls --all
  $ glab issue list --mine
  $ glab issue list --milestone release-2.0.0 --opened
  

Options
~~~~~~~

::

  -A, --all                    Get all issues
  -a, --assignee string        Filter issue by assignee <username>
      --author string          Filter issue by author <username>
  -c, --closed                 Get only closed issues
  -C, --confidential           Filter by confidential issues
      --in string              search in {title|description} (default "title,description")
  -l, --label strings          Filter issue by label <name>
  -m, --milestone string       Filter issue by milestone <id>
      --not-assignee strings   Filter issue by not being assigneed to <username>
      --not-author strings     Filter by not being by author(s) <username>
      --not-label strings      Filter issue by lack of label <name>
  -p, --page int               Page number (default 1)
  -P, --per-page int           Number of items to list per page. (default 30) (default 30)
      --search string          Search <string> in the fields defined by --in

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

