MGit: Multiple Git repository handler
=====================================

MGit is a tool to run Git commands for multiple repositories simultaneously.

MGit allows you to easily define filters and then run commands on the resulting repository list. Filters can be configured in a configuration file or you can manually enter them on the command-line.
Useful commands are installed by default. Additional commands can be added into the configuration. All existing commands can be reconfigured in the same configuration.
Results are collected and presented as one report.

Getting started
===============

Installation
------------

(work in progress)


Building
--------

git clone https://github.com/marcelfw/mgit.git
cd mgit
go build src/github.com/marcelfw/mgit


Usage examples
--------------

    # Copy your development code to your laptop.
    mgit -branch develop push laptop

    # Refresh all customer code on your machine.
    mgit -root ~/customer pull

    # Refresh all your github clones.
    mgit -remotepath github.com/username pull

    # Mirror all repositories to your NAS.
    # 1. Create a list of commands to set up remotes to push/pull to/from.
    mgit -noremote mynas echo "mkdir -p {{ .Name }}.git ; cd {{ .Name }}.git ; git init --bare" > git-create.sh
    # 2. Run git-create script in correct location on your NAS.
    scp git-create.sh git@mynas
    ssh git@mynas
    sh ./git-create.sh
    rm ./git-create.sh
    exit
    # 3. Add remote "mynas" to all repository which don't have it yet.
    mgit -noremote mynas remote add mynas "ssh://git@mynas/home/git/{{ .Name }}.git"
    # 4. Push everything.
    mgit -remote mynas push


Under the hood
--------------

* The [Go](http://golang.org) programming language.

Contributing to MGit
====================

