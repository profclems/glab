.. _glab_ssh-key_get:

glab ssh-key get
----------------

Gets a single key

Synopsis
~~~~~~~~


Returns a single SSH key specified by the ID

::

  glab ssh-key get <key-id> [flags]

Examples
~~~~~~~~

::

  # Get ssh key with ID as argument
  $ glab ssh-key get 7750633
  
  # Interactive
  $ glab ssh-key get
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

