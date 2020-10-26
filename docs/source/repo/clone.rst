.. _glab_repo_clone:

glab repo clone
---------------

Clone a Gitlab repository/project

Synopsis
~~~~~~~~


Clone supports these shorthands
- repo
- namespace/repo
- namespace/group/repo
- project ID


::

  glab repo clone <command> [flags]

Examples
~~~~~~~~

::

  $ glab repo clone profclems/glab
  $ glab repo clone https://gitlab.com/profclems/glab
  $ glab repo clone profclems/glab mydirectory  # Clones repo into mydirectory
  $ glab repo clone glab   # clones repo glab for current user 
  $ glab repo clone 4356677   # finds the project by the ID provided and clones it
  

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

