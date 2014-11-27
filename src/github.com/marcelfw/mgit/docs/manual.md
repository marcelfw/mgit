MGit: Mass Git Manual
=====================

Mass git is a tool to run commands on multiple repositories, collect their output and present it in one view.


### Filtering

Repositories can be filtered by using the filters described below. Filters can be combined as needed.
Every filter can be given on the command-line (prefix with -) or used in a shortcut ini-section.

    root         specify root directory
                 (inside shortcuts the relative root-directory
                  is taken from the location of the config file)
    depth        maximum depth to recurse directories
    name         only when text partially matches repository name

    branch       only when this branch is a branch of the repository
    nobranch     only when it is not

    tag          only when this tag is a tag of the repository
    notag        only when it is not

    remote       only when this an existing remote
    noremote     only when it is not

    remoteurl    only when text partially matches a remoteurl
    noremoteurl  only when not..

An example on the command-line would be:

    mgit -root /Users/marcel/ -branch develop list

With this shortcut defined:

    [shortcut "home"]
      root = /Users/marcel/
      branch = develop

It would be:

    mgit -s home list



### Commands

Mass git is mostly about just passing regular git commands and viewing the result is a nice view. So common git
commands are available. However you can add others commands or even add your own.
Every mass git command has a name which doesn't conflict (except for help) with an actual git command. 

A list of builtin commands:

#### List

A friendly output of all found repositories. Information includes _name_, _current branch_, _abbreviated status_,
_last commit date_ and _last commit subject_.

#### Echo

Echo lets you customise your own output of the found repositories. Uses standard Go text templating.
The following values are provided: _Path_, _Name_ and _CurrentBranch_.

    mgit echo "{{ .Name }} - {{ .Path }} - {{ .CurrentBranch }}"

Simplified "list" output.

See tips and tricks to find useful examples on how to use this.

#### Exec

Exec allows you to execute any command

#### Git commands

These are Git commands which are currently builtin. Unless the "interactive" option has been enabled, the command
will be run parallel for each found repository.
If you specific "-i" before the command, mgit will assume it has to run interactively and will not parallize it.

* status, log, commit, add
* fetch, pull, push
* remote, branch, tag


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
It is adviced to never name your custom command after a normal git command.



Usage examples
--------------

    # Show repository sizes
    mgit exec du -h -d 0

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
        cd `mgit -root ~/ -name "$1" echo "{{ .Path }}" | tail -1`
    }

Source your profile and then you can use "mcd <part-of-repository-path>" to quickly go to any repository in your home directory. Change tail to head, to return the first match instead of the last.

### Quickly view projects' git status

Add a .mgit in your repository root:

    [local]
        root = .

Now you can run "mgit status" or "mgit list" anywhere in your project directory to get the status of all repositories.

Note: this example uses the fact that a relative directory inside a configuration uses the location of the configuration file as the start directory.

