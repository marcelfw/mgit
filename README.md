MGit: Multiple Git repository handler
=====================================

MGit is a tool to run Git commands for multiple repositories simultaneously.

MGit allows you to easily define filters and then run commands on the resulting repository list. Filters can be
pre-configured in a configuration file or you can manually enter them on the command-line. There are a couple of
pre-defined commands, but they can be easily re-configured and new ones can be added with only a few lines of
configuration.

Usage examples
==============

Want to update all your code with the shared repository server?

    # Push all repositories which are on branch "develop" and have a remote "git-server"
    mgit -b develop -r git-server push git-server

    # Pull all repositories below a certain directory.
    mgit -root ~/client pull



Under the hood
--------------

* The [Go](http://golang.org) programming language.

Contributing to MGit
====================

