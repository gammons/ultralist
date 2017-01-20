# Todolist

[![](https://goreportcard.com/badge/github.com/gammons/todolist)](https://goreportcard.com/report/github.com/gammons/todolist)
[![Build Status](https://travis-ci.org/gammons/todolist.svg?branch=master)](https://travis-ci.org/gammons/todolist)

Todolist is a simple and very fast task manager for the command line.  It is based on the [Getting Things Done][gtd] methodology.

[gtd]: http://lifehacker.com/productivity-101-a-primer-to-the-getting-things-done-1551880955

## Documentation

See [The main Todolist website][tdl] for the current documentation.

[tdl]: http://todolist.site

## Is it good?

Yes.  Yes it is.

## Quick Start using [Docker](https://github.com/docker/docker.git)

Building the docker image:
```
$ git clone https://github.com/gammons/todolist.git
$ cd todolist
$ docker build -t todolist .
```

If you have an existing todo-list file, then run the container with
```
$ docker run -d -v your/todos/json/file:/.todos.json -p 7890:7890 todolist
```

Otherwise, you have to provide an empty todo file:
```
$ echo [] > /tmp/.todos.json
$ docker run -d -v /tmp/.todos.json:/.todos.json -p 7890:7890 todolist
```

Finally, open your browser and enter `localhost:7890` at the URL, and \
happy GTD!

## Author

Please send complaints, complements, rants, etc to [Grant Ammons][ga]

## License

Todolist is open source, and uses the [MIT license](https://github.com/gammons/todolist/blob/master/LICENSE.md).

[ga]: https://twitter.com/gammons
