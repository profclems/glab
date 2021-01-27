.. _glab_issue:

glab issue
----------

Work with GitLab issues

Synopsis
~~~~~~~~


Work with GitLab issues

Examples
~~~~~~~~

::

  $ glab issue list
  $ glab issue create --label --confidential
  $ glab issue view --web
  $ glab issue note -m "closing because !123 was merged" <issue number>
  

Options
~~~~~~~

::

  -R, --repo OWNER/REPO   Select another repository using the OWNER/REPO or `GROUP/NAMESPACE/REPO` format or full URL or git URL

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --help   Show help for command

Subcommands
~~~~~~~~~~~
.. toctree::
   :glob:
   :maxdepth: 0

   board <board>
   close <close>
   create <create>
   delete <delete>
   list <list>
   note <note>
   reopen <reopen>
   subscribe <subscribe>
   unsubscribe <unsubscribe>
   update <update>
   view <view>
   


