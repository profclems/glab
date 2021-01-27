.. _glab_repo_fork:

glab repo fork
--------------

Create a fork of a GitLab repository

Synopsis
~~~~~~~~


Create a fork of a GitLab repository

::

  glab repo fork <repo> [flags]

Examples
~~~~~~~~

::

  $ glab repo fork
  $ glab repo fork namespace/repo
  $ glab repo fork namespace/repo --clone
  

Options
~~~~~~~

::

  -c, --clone         Clone the fork {true|false}
  -n, --name string   The name assigned to the resultant project after forking
  -p, --path string   The path assigned to the resultant project after forking
      --remote        Add remote for fork {true|false}

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

