.. _glab_mr:

glab mr
-------

Create, view and manage merge requests

Synopsis
~~~~~~~~


Create, view and manage merge requests

Examples
~~~~~~~~

::

  $ glab mr create --autofill --labels bugfix
  $ glab mr merge 123
  $ glab mr note -m "needs to do X before it can be merged" branch-foo
  

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

   approve <approve>
   approvers <approvers>
   checkout <checkout>
   close <close>
   create <create>
   delete <delete>
   diff <diff>
   for <for>
   issues <issues>
   list <list>
   merge <merge>
   note <note>
   rebase <rebase>
   reopen <reopen>
   revoke <revoke>
   subscribe <subscribe>
   todo <todo>
   unsubscribe <unsubscribe>
   update <update>
   view <view>
   


