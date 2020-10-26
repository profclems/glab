.. _glab_alias_set:

glab alias set
--------------

Set an alias.

Synopsis
~~~~~~~~


Declare a word as a command alias that will expand to the specified command(s).

The expansion may specify additional arguments and flags. If the expansion
includes positional placeholders such as '$1', '$2', etc., any extra arguments
that follow the invocation of an alias will be inserted appropriately.

If '--shell' is specified, the alias will be run through a shell interpreter (sh). This allows you
to compose commands with "|" or redirect with ">". Note that extra arguments following the alias
will not be automatically passed to the expanded expression. To have a shell alias receive
arguments, you must explicitly accept them using "$1", "$2", etc., or "$@" to accept all of them.

Platform note: on Windows, shell aliases are executed via "sh" as installed by Git For Windows. If
you have installed git on Windows in some other way, shell aliases may not work for you.
Quotes must always be used when defining a command as in the examples.


::

  glab alias set <alias name> '<command>' [flags]

Examples
~~~~~~~~

::

  $ glab alias set mrv 'mr view'
  $ glab mrv -w 123
  #=> glab mr view -w 123
  
  $ glab alias set createissue 'glab create issue --title "$1"'
  $ glab createissue "My Issue" --description "Something is broken."
  # => glab create issue --title "My Issue" --description "Something is broken."
  
  $ glab alias set --shell igrep 'glab issue list --assignee="$1" | grep $2'
  $ glab igrep user foo
  #=> glab issue list --assignee="user" | grep "foo"
  

Options
~~~~~~~

::

  -s, --shell   Declare an alias to be passed through a shell interpreter

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

