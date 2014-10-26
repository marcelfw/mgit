MGit: Multiple Git repository handler
=====================================

MGit is a tool to run Git commands for multiple repositories simultaneously.

MGit allows you to easily define filters and then run commands on the resulting repository list. Filters can be configured in a configuration file or you can manually enter them on the command-line.
Useful commands are installed by default. Additional commands can be added into the configuration. All existing commands can be reconfigured in the same configuration.
Results are collected and presented as one report.

Usage examples
--------------

    # Copy your development code to your laptop.
    mgit -b develop push laptop

    # Refresh all customer code on your machine.
    mgit -root ~/customer pull

    # Refresh all your github clones.
    mgit -rp github.com/username pull


Under the hood
--------------

* The [Go](http://golang.org) programming language.

Contributing to MGit
====================

