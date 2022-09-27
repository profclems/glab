.. _glab_repo_view:

glab repo view
--------------

View a project/repository

Synopsis
~~~~~~~~


Display the description and README of a project or open it in the browser.

::

  glab repo view [repository] [flags]

Examples
~~~~~~~~

::

  # view project information for the current directory
  $ glab repo view
  
  # view project information of specified name
  $ glab repo view my-project
  $ glab repo view user/repo
  $ glab repo view group/namespace/repo
  
  # specify repo by full [git] URL
  $ glab repo view git@gitlab.com:user/repo.git
  $ glab repo view https://gitlab.company.org/user/repo
  $ glab repo view https://gitlab.company.org/user/repo.git
  

Options
~~~~~~~

::

  -b, --branch string   View a specific branch of the repository
  -w, --web             Open a project in the browser

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

