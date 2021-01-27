.. _glab_repo_search:

glab repo search
----------------

Search for GitLab repositories and projects by name

Synopsis
~~~~~~~~


Search for GitLab repositories and projects by name

::

  glab repo search [flags]

Examples
~~~~~~~~

::

  $ glab project search title
  $ glab repo search title
  $ glab project find title
  $ glab proejct lookup title
  

Options
~~~~~~~

::

  -p, --page int        Page number (default 1)
  -P, --per-page int    Number of items to list per page (default 20)
  -s, --search string   A string contained in the project name

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

