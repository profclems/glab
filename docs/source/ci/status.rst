.. _glab_ci_status:

glab ci status
--------------

View a running CI pipeline on current or other branch specified

Synopsis
~~~~~~~~


View a running CI pipeline on current or other branch specified

::

  glab ci status [flags]

Examples
~~~~~~~~

::

  $ glab ci status --live
  $ glab ci status --compact // more compact view
  $ glab ci status --branch=master   // Get pipeline for master branch
  $ glab ci status   // Get pipeline for current branch
  

Options
~~~~~~~

::

  -b, --branch string   Check pipeline status for a branch. (Default is current branch)
  -c, --compact         Show status in compact format
  -l, --live            Show status in real-time till pipeline ends

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

