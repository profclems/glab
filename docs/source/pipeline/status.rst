.. _glab_pipeline_status:

glab pipeline status
--------------------

View a running pipeline on current or other branch specified

Synopsis
~~~~~~~~


View a running pipeline on current or other branch specified

::

  glab pipeline status [flags]

Examples
~~~~~~~~

::

  $ glab pipeline status --live
  $ glab pipeline status --branch=master   // Get pipeline for master branch
  $ glab pipe status   // Get pipeline for current branch
  

Options
~~~~~~~

::

  -b, --branch string   Check pipeline status for a branch. (Default is current branch)
  -l, --live            Show status in real-time till pipeline ends

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

