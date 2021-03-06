MGit: Mass Git Manual
=====================

Mass git is a tool to run commands on multiple repositories, collect their output and present it in one view.


### Filtering

Repositories can be filtered by using the filters described below. Filters can be combined as needed.
Every filter can be given on the command-line (prefix with -) or used in a shortcut ini-section.

    root         specify root directory (inside shortcuts the relative root-directory
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

With this shortcut defined (see Configuration section below):

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

Exec allows you to execute any command. The working directory for the command is the actual repository directory. Just like “echo” you can use Go text templating.

Show the disk usage for all found repositories:

    mgit exec du -h -d 0

#### Git commands

These are Git commands which are currently builtin. The command
will be run parallel for each found repository.
Default Git commands available:

* status, log, commit, add
* fetch, pull, push
* remote, branch, tag

If you specific "-i" before the command, mgit will assume it has to run interactively and will not parallize it.

Run vi for each found repository:

    mgit -i exec vi .git/config


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

