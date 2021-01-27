.. _glab_ci_lint:

glab ci lint
------------

Checks if your .gitlab-ci.yml file is valid.

Synopsis
~~~~~~~~


Checks if your .gitlab-ci.yml file is valid.

::

  glab ci lint [flags]

Examples
~~~~~~~~

::

  $ glab ci lint  
  #=> Uses .gitlab-ci.yml in the current directory
  
  $ glab ci lint .gitlab-ci.yml
  
  $ glab ci lint path/to/.gitlab-ci.yml
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

