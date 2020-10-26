.. _glab_mr_list:

glab mr list
------------

List merge requests

Synopsis
~~~~~~~~


List merge requests

::

  glab mr list [flags]

Options
~~~~~~~

::

  -a, --all                Get all merge requests
      --assignee strings   Get only merge requests assigned to users
  -c, --closed             Get only closed merge requests
  -l, --label string       Filter merge request by label <name>
  -m, --merged             Get only merged merge requests
      --milestone string   Filter merge request by milestone <id>
      --mine               Get only merge requests assigned to me
  -o, --opened             Get only open merge requests
  -p, --page int           Page number (default 1)
  -P, --per-page int       Number of items to list per page. (default 30) (default 30)

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

