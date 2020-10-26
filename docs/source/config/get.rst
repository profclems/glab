.. _glab_config_get:

glab config get
---------------

Prints the value of a given configuration key

Synopsis
~~~~~~~~


Get the value for a given configuration key.

::

  glab config get <key> [flags]

Examples
~~~~~~~~

::

  
    $ glab config get editor
    vim
    $ glab config get glamour_style
    notty
  

Options
~~~~~~~

::

  -g, --global        Read from global config file (~/.config/glab-cli/config.yml). [Default: looks through Environment variables → Local → Global]
  -h, --host string   Get per-host setting

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

