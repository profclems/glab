.. _glab_ci_view:

glab ci view
------------

View, run, trace/logs, and cancel CI jobs current pipeline

Synopsis
~~~~~~~~


Supports viewing, running, tracing, and canceling jobs.

Use arrow keys to navigate jobs and logs.

'Enter' to toggle a job's logs or trace.
'Ctrl+R', 'Ctrl+P' to run/retry/play a job -- Use Tab / Arrow keys to navigate modal and Enter to confirm.
'Ctrl+C' to cancel job -- (Quits CI view if selected job isn't running or pending).
'Ctrl+Q' to Quit CI View.
'Ctrl+Space' suspend application and view logs (similar to glab pipeline ci trace)
Supports vi style (hjkl,Gg) bindings and arrow keys for navigating jobs and logs.


::

  glab ci view [branch/tag] [flags]

Examples
~~~~~~~~

::

  $ glab pipeline ci view   # Uses current branch
  $ glab pipeline ci view master  # Get latest pipeline on master branch
  $ glab pipeline ci view -b master  # just like the second example
  $ glab pipeline ci view -b master -R profclems/glab  # Get latest pipeline on master branch of profclems/glab repo
  

Options
~~~~~~~

::

  -b, --branch string   Check pipeline status for a branch/tag. (Default is the current branch)

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

