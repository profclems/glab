.. _glab_ci_delete:

glab ci delete
--------------

Delete a CI pipeline

Synopsis
~~~~~~~~


Delete a CI pipeline

::

  glab ci delete <id> [flags]

Examples
~~~~~~~~

::

  $ glab ci delete 34
  $ glab ci delete 12,34,2
  

Options
~~~~~~~

::

  -s, --status string   delete pipelines by status: {running|pending|success|failed|canceled|skipped|created|manual}

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

