.. _glab_repo_archive:

glab repo archive
-----------------

Get an archive of the repository.

Synopsis
~~~~~~~~


Clone supports these shorthands
- repo
- namespace/repo
- namespace/group/repo


::

  glab repo archive <command> [flags]

Examples
~~~~~~~~

::

  $ glab repo archive profclems/glab
  $ glab repo archive  # Downloads zip file of current repository
  $ glab repo archive profclems/glab mydirectory  # Downloads repo zip file into mydirectory
  $ glab repo archive profclems/glab --format=zip   # Finds repo for current user and download in zip format
  

Options
~~~~~~~

::

  -f, --format string   Optionally Specify format if you want a downloaded archive: {tar.gz|tar.bz2|tbz|tbz2|tb2|bz2|tar|zip} (Default: zip) (default "zip")
  -s, --sha string      The commit SHA to download. A tag, branch reference, or SHA can be used. This defaults to the tip of the default branch if not specified

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

