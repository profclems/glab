.. _glab_label_create:

glab label create
-----------------

Create labels for repository/project

Synopsis
~~~~~~~~


Create labels for repository/project

::

  glab label create [flags]

Examples
~~~~~~~~

::

  $ glab label create
  $ glab label new
  $ glab label create -R owner/repo
  

Options
~~~~~~~

::

  -c, --color string         Color of label in plain or HEX code. (Default: #428BCA) (default "#428BCA")
  -d, --description string   Label description
  -n, --name string          Name of label

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

