.. _glab_repo_mirror:

glab repo mirror
----------------

Mirror a project/repository

Synopsis
~~~~~~~~


Mirrors a project/repository to the specified location using pull or push method.

::

  glab repo mirror [ID | URL | PATH] [flags]

Options
~~~~~~~

::

      --allow-divergence          Determines if divergent refs are skipped.
      --direction string          Mirror direction (default "pull")
      --enabled                   Determines if the mirror is enabled. (default true)
      --protected-branches-only   Determines if only protected branches are mirrored.
      --url string                The target URL to which the repository is mirrored.

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

