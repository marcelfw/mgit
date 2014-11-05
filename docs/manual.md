MGit: Mass Git Manual
=====================

Mass git is a tool to run commands on multiple repositories, collect their output and present it in one view.


Features detailed
-----------------

### Filtering

Repositories can be filtered by using the filters as described below. Filters can be combined as necessary.

#### Root directory and recurse depth

The root directory specifies the start of the repository search and the recurse depth limits how deep it will search.

    # Search my home directory for 2 levels and outputs all repositories it finds.
    mgit -root /Users/marcel -depth 2 list

If not otherwise specified the default -root is the current directory. Depth is unlimited if not specified.

#### Branch and remote

You can filter on the presence of branch or remote, or on the absence of one.

    # List repositories which have a 'develop' branch.
    mgit -branch develop list

    # List repositories which don't have a remote called 'laptop'.
    mgit -noremote laptop list

#### Remote path

Sometimes you don't know (or don't care) about the name of the remote, but you do know which path you would like
or not.

    # List my repositories cloned from github.
    mgit -remotepath github.com/marcelfw list

#### Name

The name of the repository is the full path excluding the root directory. You can search for a specific repository
using this.

    # List repositories in a sub-directory called 'client-42'.
    mgit -name /client-42/ list


Commands
--------

MGit only contains a few builtin commands, the rest are calls to Git itself. However the builtins can be quite useful.


#### List

A friendly output of all found repositories. Information includes _name_, _current branch_, _abbreviated status_,
_last commit date_ and _last commit subject_.

#### Echo

Echo lets you customize your own output of the found repositories. Uses standard Go text templating.
The following values are provided: _Path_, _Name_ and _CurrentBranch_.

    mgit echo "{{ .Name }} - {{ .Path }} - {{ .CurrentBranch }}"

Simplified "list" output.

See tips and tricks to find useful examples on how to use this.

#### Git commands

These are Git commands which are currently builtin. Unless the "interactive" option has been enabled, the command
will be run parallel for each found repository.
If you specific "-i" before the command, mgit will assume it has to run interactively and will not parallize it.

* status, log, commit, add
* fetch, pull, push
* remote, branch


Configuration
-------------

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
    # 2. Add remote "mynas" to all repositories which don not have it yet.
    mgit -noremote mynas remote add mynas "ssh://git@mynas/home/git/{{ .Name }}.git"
    # 3. Push everything.
    mgit -remote mynas push





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
