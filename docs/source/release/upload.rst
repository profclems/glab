.. _glab_release_upload:

glab release upload
-------------------

Upload release asset files or links to GitLab Release

Synopsis
~~~~~~~~


Upload release assets to GitLab Release

You can define the display name by appending '#' after the file name. 
The link type comes after the display name (eg. 'myfile.tar.gz#My display name#package')


::

  glab release upload <tag> [<files>...] [flags]

Examples
~~~~~~~~

::

  Upload a release asset with a display name
  $ glab release upload v1.0.1 '/path/to/asset.zip#My display label'
  
  Upload a release asset with a display name and type
  $ glab release upload v1.0.1 '/path/to/asset.png#My display label#image'
  
  Upload all assets in a specified folder
  $ glab release upload v1.0.1 ./dist/*
  
  Upload all tarballs in a specified folder
  $ glab release upload v1.0.1 ./dist/*.tar.gz
  
  Upload release assets links specified as JSON string
  $ glab release upload v1.0.1 --assets-links='
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

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

