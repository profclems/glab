.. _glab_repo_create:

glab repo create
----------------

Create a new GitLab project/repository

Synopsis
~~~~~~~~


Create a new GitLab repository.

::

  glab repo create [path] [flags]

Examples
~~~~~~~~

::

  # create a repository under your account using the current directory name
  $ glab repo create
  
  # create a repository under a group using the current directory name
  $ glab repo create --group glab-cli
  
  # create a repository with a specific name
  $ glab repo create my-project
  
  # create a repository for a group
  $ glab repo create glab-cli/my-project
  

Options
~~~~~~~

::

      --defaultBranch master   Default branch of the project. If not provided, master by default.
  -d, --description string     Description of the new project
  -g, --group string           Namespace/group for the new project (defaults to the current user’s namespace)
      --internal               Make project internal: visible to any authenticated user (default)
  -n, --name string            Name of the new project
  -p, --private                Make project private: visible only to project members
  -P, --public                 Make project public: visible without any authentication
      --readme                 Initialize project with README.md
      --remoteName origin      Remote name for the Git repository you're in. If not provided, origin by default. (default "origin")
  -t, --tag stringArray        The list of tags for the project.

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

