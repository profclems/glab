.. _glab_ci_list:

glab ci list
------------

Get the list of CI pipelines

Synopsis
~~~~~~~~


Get the list of CI pipelines

::

  glab ci list [flags]

Examples
~~~~~~~~

::

  $ glab ci list
  $ glab ci list --state=failed
  

Options
~~~~~~~

::

  -o, --orderBy string   Order pipeline by <string>
  -p, --page int         Page number (default 1)
  -P, --per-page int     Number of items to list per page. (default 30) (default 30)
      --sort string      Sort pipeline by {asc|desc}. (Defaults to desc) (default "desc")
  -s, --status string    Get pipeline with status: {running|pending|success|failed|canceled|skipped|created|manual}

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

