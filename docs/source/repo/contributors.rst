.. _glab_repo_contributors:

glab repo contributors
----------------------

Get contributors of the repository.

Synopsis
~~~~~~~~


Clone supports these shorthands
- repo
- namespace/repo
- namespace/group/repo


::

  glab repo contributors [flags]

Examples
~~~~~~~~

::

  $ glab repo contributors
  $ glab repo archive  // Downloads zip file of current repository
  

Options
~~~~~~~

::

  -f, --order string      Return contributors ordered by name, email, or commits (orders by commit date) fields. Default is commits (default "zip")
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL
  -s, --sort string       Return contributors sorted in asc or desc order. Default is asc

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

