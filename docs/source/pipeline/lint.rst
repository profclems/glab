.. _glab_pipeline_ci_lint:

glab pipeline ci lint
---------------------

Checks if your .gitlab-ci.yml file is valid.

Synopsis
~~~~~~~~


Checks if your .gitlab-ci.yml file is valid.

::

  glab pipeline ci lint [flags]

Examples
~~~~~~~~

::

  $ glab pipeline ci lint  # Uses .gitlab-ci.yml in the current directory
  $ glab pipeline ci lint .gitlab-ci.yml
  $ glab pipeline ci lint path/to/.gitlab-ci.yml
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL

