.. _glab_variable_update:

glab variable update
--------------------

Update an existing project or group variable

Synopsis
~~~~~~~~


Update an existing project or group variable

::

  glab variable update <key> <value> [flags]

Examples
~~~~~~~~

::

  $ glab variable update WITH_ARG "some value"
  $ glab variable update FROM_FLAG -v "some value"
  $ glab variable update FROM_ENV_WITH_ARG "${ENV_VAR}"
  $ glab variable update FROM_ENV_WITH_FLAG -v"${ENV_VAR}"
  $ glab variable update FROM_FILE < secret.txt
  $ cat file.txt | glab variable update SERVER_TOKEN
  $ cat token.txt | glab variable update GROUP_TOKEN -g mygroup --scope=prod
  

Options
~~~~~~~

::

  -g, --group string   Set variable for a group
  -m, --masked         Whether the variable is masked
  -p, --protected      Whether the variable is protected
  -s, --scope string   The environment_scope of the variable. All (*), or specific environments (default "*")
  -t, --type string    The type of a variable: {env_var|file} (default "env_var")
  -v, --value string   The value of a variable

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

