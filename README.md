# Todolist

Todolist is a simple and very fast task manager for the command line.  It is based on the [Getting Things Done][gtd] methodology.

## Installation

Grab the binary that makes sense for your platform:

## Usage

Todolist looks for a `.todos.json` file in the current directory you are in.

The first command you will want to run is `todo init`.

```
~ todo init
Todo repo initialized.
```

#### Adding todos

Adding todos is as simple as using the `add` or `a` command.

```
~ todo add Talk with @bob about the +devopsProject due tuesday
Todo added.
```

**Projects and Contexts**

When you add a todo, you can add one or more `+project`s or `@context`s, which you can use for grouping or filtering.

```
~ todo a Did @mary talk with @bob about the +devopsProject? due tod
Todo added.
```

**Dates**

You can also add a due date by adding `due <date>` at the end.  Below is a list of phrases that `due` can recognize:

* `due today` or `due tod`
* `due tomorrow` or `due tom`

* `due (monday|tuesday|wednesday|thursday|friday|saturday|sunday)`
* `due (mon|tue|wed|thur|fri|sat|sun)`

Sets the due date to be a specific day of this or next week, depending on what day it is.

For example, if today is a Tuesday:

* `todo a chat with @nick due wed` will be the very next day.
* `todo a chat with @nick due mon` will be the *following* Monday.

You can also use specific dates, provided the are in the format of `jan 2` or `2 jan`.

* `todo a chat with @nick due may 15`
* `todo a chat with @nick due 15 may`

#### Listing todos

## FAQ

**Really? Another task manager?**

I wrote todolist because I basically wanted to be able to grep my todos easily.  Also, while Wunderlist is fantastic, managing todos is something that is easily done on the command line.
