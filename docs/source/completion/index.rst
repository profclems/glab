.. _glab_completion:

glab completion
---------------

Generate shell completion scripts

Synopsis
~~~~~~~~


Generate shell completion scripts for glab commands.

The output of this command will be computer code and is meant to be saved to a
file or immediately evaluated by an interactive shell.

For example, for bash you could add this to your '~/.bash_profile':

	eval "$(glab completion -s bash)"

When installing glab through a package manager, however, it's possible that
no additional shell configuration is necessary to gain completion support. 
For Homebrew, see <https://docs.brew.sh/Shell-Completion>


::

  glab completion [flags]

Options
~~~~~~~

::

  -s, --shell string   Shell type: {bash|zsh|fish|powershell}

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

