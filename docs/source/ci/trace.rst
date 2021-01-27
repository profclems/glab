.. _glab_ci_trace:

glab ci trace
-------------

Trace a CI job log in real time

Synopsis
~~~~~~~~


Trace a CI job log in real time

::

  glab ci trace [<job-id>] [flags]

Examples
~~~~~~~~

::

  $ glab ci trace
  #=> interactively select a job to trace
  
  $ glab ci trace 224356863
  #=> trace job with id 224356863
  

Options
~~~~~~~

::

  -b, --branch string   Check pipeline status for a branch. (Default is the current branch)

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

