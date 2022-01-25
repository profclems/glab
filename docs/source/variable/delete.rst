.. _glab_variable_delete:

glab variable delete
--------------------

Delete a project or group variable

Synopsis
~~~~~~~~


Delete a project or group variable

::

  glab variable delete <key> [flags]

Examples
~~~~~~~~

::

  $ glab variable delete VAR_NAME
  $ glab variable delete VAR_NAME --scope=prod
  $ glab variable delete VARNAME -g mygroup
  

Options
~~~~~~~

::

  -g, --group string   Delete variable from a group
  -s, --scope string   The environment_scope of the variable. All (*), or specific environments (default "*")

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

