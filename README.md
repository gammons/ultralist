# Ultralist
### Simple task management for tech folks.

[![](https://goreportcard.com/badge/github.com/ultralist/ultralist)](https://goreportcard.com/report/github.com/ultralist/ultralist)
[![Actions Status](https://github.com/ultralist/ultralist/workflows/Go/badge.svg)](https://github.com/ultralist/ultralist/actions)

Ultralist is a task management system for technical people. It is command-line component that is very fast and stays out of the way. 

[![](https://ultralist.io/images/ultralist.png)](https://ultralist.io)

Ultralist is based off of the [Getting Things Done][gtd] system, and is centered around the following concepts:

* due dates
* projects and contexts
* statuses
* task recurrence


The CLI is _fast_, _powerful_ and _intuitive_.  It will also always be open source.

## Documentation

All of Ultralist's documentation is available on the [Ultralist website](https://ultralist.io).

* [Ultralist Concepts](https://ultralist.io/docs/basics/concepts/)
* [Quickstart](https://ultralist.io/docs/cli/quickstart/)
* [Managing todos](https://ultralist.io/docs/cli/managing_tasks/)
* [Todo Recurrence](https://ultralist.io/docs/cli/recurrence/)
* [Listing and filtering todos](https://ultralist.io/docs/cli/showing_tasks/)
* [Best Practices](https://ultralist.io/docs/cli/best_practices/)
* [Syncing with Ultralist Pro](https://ultralist.io/docs/cli/pro_integration/)
* [The .todos.json file format](https://ultralist.io/docs/cli/todos_json/)

## Ultralist Pro

You can optionally combine the Ultralist CLI with [Ultralist Pro](https://ultralist.io).  Doing so adds the following benefits:

* Easily keep CLI lists in sync across multiple computers.
* Manage your list with a slick web app.
* Use the Ultralist mobile apps.
* Use the Slack integration. Add + manage tasks directly from Slack.
* Use our robust API to enable more complex workflows.

Ultralist Pro provides a superior task management experience to Todoist, Any.do etc.  The command-line will app _always_ be first and foremost.

## Is it good?

Yes.  Yes it is.

## Installation

* **Mac OS**: Run `brew install ultralist`. (Or `port install ultralist` if you are using [MacPorts](https://www.macports.org).)
* **Arch Linux**: May be installed from AUR [ultralist](https://aur.archlinux.org/packages/ultralist/)
* **FreeBSD**: Run `pkg install ultralist` or `cd /usr/ports/deskutils/ultralist && make install clean`
* **Other systems**: Get the correct ultralist binary from the [releases page](https://github.com/ultralist/ultralist/releases).
* If you have Golang installed: Run `go get github.com/ultralist/ultralist`.

Then, follow the [quick start](https://ultralist.io/docs/cli/quickstart/) in the docs to quickly get up and running.

## How is this different from todo.txt, Taskwarrior, etc?

[todo.txt](http://todotxt.org/) is great.  But it didn't work well for my needs.

Here's how ultralist compares specifically to todo.txt:
1. **Due dates.** they are a core concept of ultralist but not todo.txt.
1. **Synchronizing.** Syncing is built into the CLI using the [ultralist.io](https://ultralist.io) service.
1. **Active development.** the ultralist CLI is under active development, whereas todo.txt's CLI is not.

Taskwarrior is a similar system, however it is less intuitive and not maintained.

## Author

Please send complaints, complements, rants, etc to [Grant Ammons][ga]

## License

Ultralist is open source, and uses the [MIT license](https://github.com/ultralist/ultralist/blob/master/LICENSE.md).

[ga]: https://twitter.com/gammons
[gtd]: http://lifehacker.com/productivity-101-a-primer-to-the-getting-things-done-1551880955
