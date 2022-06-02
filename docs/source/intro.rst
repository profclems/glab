GLab - A GitLab CLI Tool
------------------------
GLab is an open source Gitlab Cli tool written in Go (golang) to help
work seamlessly with Gitlab from the command line. Work with issues,
merge requests, **watch running pipelines directly from your CLI** among
other features. Inspired by ``gh``, `the official GitHub CLI
tool <https://github.com/cli/cli>`__.

Usage
-----
.. code:: sh

   glab <command> <subcommand> [flags]

Core Commands
~~~~~~~~~~~~~

-  ``glab mr [list, create, close, reopen, delete, ...]``
-  ``glab issue [list, create, close, reopen, delete, ...]``
-  ``glab pipeline [list, delete, status, view, ...]``
-  ``glab release``
-  ``glab repo``
-  ``glab label``
-  ``glab alias``

Examples
~~~~~~~~
.. code:: sh

    $ glab auth login --stdin < token.txt
    $ glab issue list
    $ glab mr for 123   # Create merge request for issue 123
    $ glab mr checkout 243
    $ glab pipeline ci view
    $ glab mr view
    $ glab mr approve
    $ glab mr merge

Installation
~~~~~~~~~~~~
You can find installation instructions on our `README <https://github.com/profclems/glab#installation>`__.

Authentication
~~~~~~~~~~~~~~
Run ``glab auth login`` to authenticate with your GitLab account. ``glab`` will respect tokens set using ``GITLAB_TOKEN``.

Feedback
~~~~~~~~
Thank you for checking out GLab! Please open an `issue <https://github.com/profclems/glab/issues/new>`__. to send us feedback. We're looking forward to hearing it.
