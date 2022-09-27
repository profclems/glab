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

Generate a %[1]s_gh%[1]s completion script and put it somewhere in your %[1]s$fpath%[1]s:
				gh completion -s zsh > /usr/local/share/zsh/site-functions/_gh
			Ensure that the following is present in your %[1]s~/.zshrc%[1]s:
				autoload -U compinit
				compinit -i
			
			Zsh version 5.7 or later is recommended.

When installing glab through a package manager, however, it's possible that
no additional shell configuration is necessary to gain completion support. 
For Homebrew, see <https://docs.brew.sh/Shell-Completion>


::

  glab completion [flags]

Options
~~~~~~~

::

      --no-desc        Do not include shell completion description
  -s, --shell string   Shell type: {bash|zsh|fish|powershell} (default "bash")

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

