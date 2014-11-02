MGit: Multiple Git repository handler
=====================================

MGit is a tool to run Git commands for multiple repositories simultaneously.

MGit allows you to easily define filters and then run commands on the resulting repository list. Filters can be
configured in a file or you can manually enter them on the command-line.
Useful commands are installed by default. Additional commands can be added into the configuration. All existing
commands can be reconfigured in the same configuration.
Results are collected and presented as one report.

Getting started
===============

There is no binary release yet.


Building from source
--------------------

git clone https://github.com/marcelfw/mgit.git
cd mgit
go build src/github.com/marcelfw/mgit


Usage examples
--------------

    # Copy all your development code branches to your laptop.
    mgit -branch develop -remote laptop push laptop

    # Refresh all customer code on your machine.
    mgit -root ~/customer pull

    # Refresh all your github clones.
    mgit -remotepath github.com/username pull

    # Mirror all repositories to your NAS.
    # 1. Create script and feed into NAS with ssh (shell should allow for git init).
    mgit -noremote mynas echo "mkdir -p {{ .Name }}.git ; git init --bare {{ .Name }}.git" | ssh git@mynas
    # 2. Add remote "mynas" to all repositories which don't have it yet.
    mgit -noremote mynas remote add mynas "ssh://git@mynas/home/git/{{ .Name }}.git"
    # 3. Push everything.
    mgit -remote mynas push


Under the hood
--------------

* The [Go](http://golang.org) programming language.

Contributing to MGit
====================

