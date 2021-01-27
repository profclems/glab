.. _glab_variable_set:

glab variable set
-----------------

Create a new project or group variable

Synopsis
~~~~~~~~


Create a new project or group variable

::

  glab variable set <key> <value> [flags]

Examples
~~~~~~~~

::

  $ glab variable set WITH_ARG "some value"
  $ glab variable set FROM_FLAG -v "some value"
  $ glab variable set FROM_ENV_WITH_ARG "${ENV_VAR}"
  $ glab variable set FROM_ENV_WITH_FLAG -v"${ENV_VAR}"
  $ glab variable set FROM_FILE < secret.txt
  $ cat file.txt | glab variable set SERVER_TOKEN
  $ cat token.txt | glab variable set GROUP_TOKEN -g mygroup --scope=prod
  

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

