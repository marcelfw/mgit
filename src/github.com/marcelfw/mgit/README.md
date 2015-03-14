MGit: Mass Git
==============

Mass git is a tool to run commands on multiple repositories, collect their output and present it in one view.

* Filter repositories on branch, remote, name and more
* Execute Git commands on multiple repositories
* Combine all outputs
* Simple configuration
* Customize and extend commands


Examples
--------

    # Get an overview of all customer code on your machine.
    mgit -root ~/customer list

    # Alternatively use the regular git status command.
    mgit -root ~/customer status

    # Copy all your development branches to your laptop.
    mgit -branch develop -remote laptop push laptop

    # Refresh your github clones.
    mgit -remotepath github.com/username pull

    # Mirror all repositories to your NAS.
    # 1. Create bare repositories onto your NAS with ssh (shell should allow for git init).
    mgit -noremote mynas exec ssh git@mynas mkdir -p "{{ .Name }}" \; git init --bare "{{ .Name }}"
    # 2. Add remote "mynas" to all repositories which don't have it yet.
    mgit -noremote mynas remote add mynas "ssh://git@mynas/home/git/{{ .Name }}.git"
    # 3. Push everything.
    mgit -remote mynas push

    # Go to any repository directory, in this case project "mgit"
    # (see Tips 'n tricks)
    mcd mgit


Getting started
---------------

There is no binary release yet. See next section.


Installing from source
----------------------

1. Install Go, see Go [documentation](http://golang.org/doc/install)
3. go get github.com/marcelfw/mgit
4. Optionally copy bin/mgit to /usr/local/bin directory.


Tips 'n tricks
--------------

### Quickly go to directory of a repository

Add a shortcut called "global" into your system or user-global configuration. Set the "root" to your
root of all repositories.
Source the following snippet in your shell profile:

source &lt;mgit-directory&gt;/profile-mgit.sh

Source your profile to load the changes immediately.
Now you can use "mcd" in two ways:

1. Use `mcd <name>` to go to the repository which matches this <name>
2. Or `mcd .` to go to the local mgit root

### Quickly view projects' git status

Add a .mgit in your repository root:

    [local]
        root = .

Now you can run "mgit status" or "mgit list" anywhere in your project directory to get the status of all repositories.


License
-------

Code is under the [The MIT License (MIT)](https://github.com/marcelfw/mgit/blob/master/LICENSE.txt).
