# dom ( digital-ocean-manager )
A command line tool for interacting with your DigitalOcean droplets.

This tool in inspired by [TugBoat](https://github.com/pearkes/tugboat) tool, which is implemented in ruby.
This is go implementation of same.

NOTE
=====
This is still in development and all features are yet not available.
Current Status
---------------
NOT READY. Under active development.


Installation.
-------------
- Still no release is done.
- This is go getable project, execute
   - `` $ go get github.com/kunalkushwaha/dom``


prerequieste
------------
golang 1.4+

Usage.
-------
```
$ dom help
NAME:
   digital-ocean-manager - Cli based digital-ocean client

USAGE:
   digital-ocean-manager [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   images	List images
   create	creates a droplet
   destory	destroy a droplet
   droplets	list  droplets
   regions	list a regions
   resize	resize a droplet
   halt		shutdown droplet
   restart	restart droplet
   info		details of droplet
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version

```

Reporting Bugs
----------------
Yes, please!
You can create a new issue here. Thank you!

Contributing
-------------
1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request


Development Environment
------------------------
To add a feature, fix a bug, or to run a development build of dom on your machine.

```
$ go get github.com/tools/godep
$ go get github.com/<your-github-id> dom
$ cd $GOPATH/src/github.com/<your-github-id/dom
$ godep restore
$ go build
```

To run test cases
```
$ go test
```

If you need help with your environment, feel free to open an issue.
