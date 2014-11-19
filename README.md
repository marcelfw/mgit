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

    # Copy all your development branches to your laptop.
    mgit -branch develop -remote laptop push laptop

    # Refresh all customer code on your machine.
    mgit -root ~/customer pull

    # Refresh all your github clones.
    mgit -remotepath github.com/username pull

    # Mirror all repositories to your NAS.
    # 1. Create script and feed into NAS with ssh (shell should allow for git init).
    mgit -noremote mynas echo "mkdir -p {{ .Name }}.git ; git init --bare {{ .Name }}.git" | ssh git@mynas
    mgit -debug -s tt exec ssh git@192.168.2.100 mkdir -p "{{ .Name }}" \; git init --bare "{{ .Name }}"
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
