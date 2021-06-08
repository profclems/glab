.. _glab_release_create:

glab release create
-------------------

Create a new or update a GitLab Release for a repository

Synopsis
~~~~~~~~


Create a new or update a GitLab Release for a repository.

If the release already exists, glab updates the release with the new info provided.

If a git tag specified does not yet exist, the release will automatically get created
from the latest state of the default branch and tagged with the specified tag name.
Use `--ref` to override this.
The `ref` can be a commit SHA, another tag name, or a branch name.
To fetch the new tag locally after the release, do `git fetch --tags origin`.

To create a release from an annotated git tag, first create one locally with
git, push the tag to GitLab, then run this command.

NB: Developer level access to the project is required to create a release.


::

  glab release create <tag> [<files>...] [flags]

Examples
~~~~~~~~

::

  Interactively create a release
  $ glab release create v1.0.1
  
  Non-interactively create a release by specifying a note
  $ glab release create v1.0.1 --notes "bugfix release"
  
  Use release notes from a file
  $ glab release create v1.0.1 -F changelog.md
  
  Upload a release asset with a display name
  $ glab release create v1.0.1 '/path/to/asset.zip#My display label'
  
  Upload a release asset with a display name and type
  $ glab release create v1.0.1 '/path/to/asset.png#My display label#image'
  
  Upload all assets in a specified folder
  $ glab release create v1.0.1 ./dist/*
  
  Upload all tarballs in a specified folder
  $ glab release create v1.0.1 ./dist/*.tar.gz
  
  Create a release with assets specified as JSON object
  $ glab release create v1.0.1 --assets-links='
  	[
  		{
  			"name": "Asset1", 
  			"url":"https://<domain>/some/location/1", 
  			"link_type": "other", 
  			"filepath": "path/to/file"
  		}
  	]'
  

Options
~~~~~~~

::

  -a, --assets-links JSON   JSON string representation of assets links (e.g. `--assets='[{"name": "Asset1", "url":"https://<domain>/some/location/1", "link_type": "other", "filepath": "path/to/file"}]')`
  -m, --milestone strings   The title of each milestone the release is associated with
  -n, --name string         The release name or title
  -N, --notes string        The release notes/description. You can use Markdown
  -F, --notes-file file     Read release notes file. Specify `-` as value to read from stdin
  -r, --ref string          If a tag specified doesn't exist, the release is created from ref and tagged with the specified tag name. It can be a commit SHA, another tag name, or a branch name.
  -D, --released-at date    The date when the release is/was ready. Defaults to the current datetime. Expected in ISO 8601 format (2019-03-15T08:00:00Z)

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

