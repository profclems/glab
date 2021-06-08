.. _glab_release_view:

glab release view
-----------------

View information about a GitLab Release

Synopsis
~~~~~~~~


View information about a GitLab Release.

Without an explicit tag name argument, the latest release in the project is shown.
%!(EXTRA string=`)

::

  glab release view <tag> [flags]

Examples
~~~~~~~~

::

  View the latest release of a GitLab repository
  $ glab release view
  
  View a release with specified tag name
  $ glab release view v1.0.1 
  

Options
~~~~~~~

::

  -w, --web   Open the release in the browser

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

