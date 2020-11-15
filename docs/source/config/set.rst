.. _glab_config_set:

glab config set
---------------

Updates configuration with the value of a given key

Synopsis
~~~~~~~~


Update the configuration by setting a key to a value.
Use glab config set --global if you want to set a global config. 
Specifying the --hostname flag also saves in the global config file


::

  glab config set <key> <value> [flags]

Examples
~~~~~~~~

::

  
    $ glab config set editor vim
    $ glab config set token xxxxx -h gitlab.com
  

Options
~~~~~~~

::

  -g, --global        write to global ~/.config/glab-cli/config.yml file rather than the repository .glab-cli/config/config
  -h, --host string   Set per-host setting

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

