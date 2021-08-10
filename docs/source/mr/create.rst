.. _glab_mr_create:

glab mr create
--------------

Create new merge request

Synopsis
~~~~~~~~


Create new merge request

::

  glab mr create [flags]

Examples
~~~~~~~~

::

  $ glab mr new
  $ glab mr create -a username -t "fix annoying bug"
  $ glab mr create -f --draft --label RFC
  $ glab mr create --fill --yes --web
  

Options
~~~~~~~

::

      --allow-collaboration    Allow commits from other members
  -a, --assignee usernames     Assign merge request to people by their usernames
      --create-source-branch   Create source branch if it does not exist
  -d, --description string     Supply a description for merge request
      --draft                  Mark merge request as a draft
  -f, --fill                   Do not prompt for title/description and just use commit info
  -H, --head OWNER/REPO        Select another head repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or the project ID or full URL
  -l, --label strings          Add label by name. Multiple labels should be comma separated
  -m, --milestone string       The global ID or title of a milestone to assign
      --no-editor              Don't open editor to enter description. If set to true, uses prompt. Default is false
      --push                   Push committed changes after creating merge request. Make sure you have committed changes
      --remove-source-branch   Remove Source Branch on merge
  -s, --source-branch string   The Branch you are creating the merge request. Default is the current branch.
  -b, --target-branch string   The target or base branch into which you want your code merged
  -t, --title string           Supply a title for merge request
  -w, --web                    continue merge request creation on web browser
      --wip                    Mark merge request as a work in progress. Alternative to --draft
  -y, --yes                    Skip submission confirmation prompt, with --fill it skips all optional prompts

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help              Show help for command
  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

