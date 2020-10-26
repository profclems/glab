.. _glab_config:

glab config
-----------

Set and get glab settings

Synopsis
~~~~~~~~


Get and set key/value strings.

Current respected settings:

- token: Your gitlab access token, defaults to environment variables
- gitlab_uri: if unset, defaults to https://gitlab.com
- browser: if unset, defaults to environment variables
- editor: if unset, defaults to environment variables.
- visual: alternative for editor. if unset, defaults to environment variables.
- glamour_style: Your desired markdown renderer style. Options are dark, light, notty. Custom styles are allowed set a custom style https://github.com/charmbracelet/glamour#styles
	

Options
~~~~~~~

::

  -g, --global   use global config file

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

Subcommands
~~~~~~~~~~~
.. toctree::
   :glob:
   :maxdepth: 0

   get <get>
   init <init>
   set <set>
   


