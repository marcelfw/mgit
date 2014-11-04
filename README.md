MGit: Multiple Git repository handler
=====================================

MGit is a tool to run Git commands on multiple repositories, collect output and present it in one view.

Features
--------

1. Filter repositories
2. Pre-configured commands
3. Combine all outputs
4. Store configuration for re-use
5. Customize and extend commands

_(see feature details below)_


Usage examples
--------------

    # Copy all your development branches to your laptop.
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


Getting started
---------------

There is no binary release yet.


Building from source
--------------------

1. Install Go, see Go [documentation](http://golang.org/doc/install)
2. git clone https://github.com/marcelfw/mgit.git
3. go build src/github.com/marcelfw/mgit/mgit.go
4. Copy binary to bin directory.


Features detailed
-----------------

#### Filter repositories

Currently you can filter on these things:

* directory and recurse depth
* branch or nobranch
* remote or noremote
* remotepath

#### Pre-configured commands

Builtin commands:

* list - list repositories
* path - show complete path
* echo - echo repository information
* help - show help information (also help "command")
* version - show version

Git "proxied" commands:

* status, log, commit, add
* fetch, pull, push
* remote, branch

By default commands will be run simultaneously but you can add an option to run them interactively.

#### Combine all outputs

If possible commands are run simultaneously and their output is collected and returned in one view.

#### Store configuration for re-use

Shortcuts allow you to combine filters and re-use them from the command-line.
Configuration are simple text files in ini-format.

Global configurations are always read and allow you to store system-global, user-global shortcuts and custom commands.

    ~/.mgit        user configuration
    /etc/mgit      system configuration

Directory configurations are searched from the current directory all the way to the root and allow you to set project defaults, shortcuts
and commands.

Directory configurations are searched first and then user- and system-configurations. The first match for a shortcut or command will be used.


#### Customize and extend commands

Pre-configured commands can be overridden in your own configuration file and you can add your own Git commands.




Tips 'n tricks
--------------

### Quickly go to directory of a repository

Add the following code to your shell profile:

    function mcd()
    {
        cd `mgit -root ~/ path $1 | tail -1`
    }

Source your profile and then you can use "mcd <part-of-repository-path>" to quickly go to any repository in your home directory.

### Quickly view projects' git status

Add a .mgit in your repository root:

    [local]
        root = .

Now you can run "mgit status" or "mgit list" anywhere in your project directory to get the status of all repositories.


License
-------

Code is under the [The MIT License (MIT)](https://github.com/marcelfw/mgit/tree/master/LICENSE.txt).