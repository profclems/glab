.. _glab_ssh-key_add:

glab ssh-key add
----------------

Add an SSH key to your GitLab account

Synopsis
~~~~~~~~


Creates a new SSH key owned by the currently authenticated user.

The --title flag is always required


::

  glab ssh-key add [key-file] [flags]

Examples
~~~~~~~~

::

  # Read ssh key from stdin and upload
  $ glab ssh-key add -t "my title"
  
  # Read ssh key from specified key file and upload
  $ glab ssh-key add ~/.ssh/id_ed25519.pub -t "my title"
  

Options
~~~~~~~

::

  -e, --expires-at string   The expiration date of the SSH key in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
  -t, --title string        New SSH key's title

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

