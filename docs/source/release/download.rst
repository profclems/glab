.. _glab_release_download:

glab release download
---------------------

Download asset files from a GitLab Release

Synopsis
~~~~~~~~


Download asset files from a GitLab Release

If no tag is specified, assets are downloaded from the latest release.
Use `--asset-name` to specify a file name to download from the release assets.
`--asset-name` flag accepts glob patterns.

Unless `--include-external` flag is specified, external files are not downloaded.


::

  glab release download <tag> [flags]

Examples
~~~~~~~~

::

  Download all assets from the latest release
  $ glab release download
  
  Download all assets from the specified release tag
  $ glab release download v1.1.0
  
  Download assets with names matching the glob pattern
  $ glab release download v1.10.1 --asset-name="*.tar.gz"
  

Options
~~~~~~~

::

  -n, --asset-name stringArray   Download only assets that match the name or a glob pattern
  -D, --dir string               Directory to download the release assets to (default ".")
  -x, --include-external         Include external asset files

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

