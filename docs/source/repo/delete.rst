.. _glab_repo_delete:

glab repo delete
----------------

Delete an existing repository on GitLab

Synopsis
~~~~~~~~


Delete an existing repository on GitLab

::

  glab repo delete [<NAMESPACE>/]<NAME> [flags]

Examples
~~~~~~~~

::

  # delete a personal repo
  $ glab repo delete dotfiles
  
  # delete a repo in GitLab group or another repo you have write access
  $ glab repo delete mygroup/dotfiles
  
  $ glab repo delete myorg/mynamespace/dotfiles
   

Options
~~~~~~~

::

  -y, --yes   Skip the confirmation prompt and immediately delete the repository.

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

