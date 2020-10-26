.. _glab_pipeline_delete:

glab pipeline delete
--------------------

Delete a pipeline

Synopsis
~~~~~~~~


Delete a pipeline

::

  glab pipeline delete <id> [flags]

Examples
~~~~~~~~

::

  $ glab pipeline delete 34
  $ glab pipeline delete 12,34,2
  

Options
~~~~~~~

::

  -s, --status string   delete pipelines by status: {running|pending|success|failed|canceled|skipped|created|manual}

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

