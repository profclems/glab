.. _glab_mr_view:

glab mr view
------------

Display the title, body, and other information about a merge request.

Synopsis
~~~~~~~~


Display the title, body, and other information about a merge request.

::

  glab mr view {<id> | <branch>} [flags]

Options
~~~~~~~

::

  -c, --comments       Show mr comments and activities
  -p, --page int       Page number (default 1)
  -P, --per-page int   Number of items to list per page (default 20)
  -s, --system-logs    Show system activities / logs
  -w, --web            Open mr in a browser. Uses default browser or browser specified in BROWSER variable

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

