.. _glab_ci_run:

glab ci run
-----------

Create or run a new CI pipeline

Synopsis
~~~~~~~~


Create or run a new CI pipeline

::

  glab ci run [flags]

Examples
~~~~~~~~

::

  $ glab ci run
  $ glab ci run -b trunk
  

Options
~~~~~~~

::

  -b, --branch string       Create pipeline on branch/ref <string>
      --variables strings   Pass variables to pipeline

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

