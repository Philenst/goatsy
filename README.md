# goatsy

*A printing tool like cutesy, but written in Go!*
Changes are expected.

If you need a JavaScript version (the original), check out [cutesy](https://github.com/calemy/cutesy).

---

##  How to use it

### Download the module

```bash
go get -u github.com/calemy/goatsy
```

### Create a logger

I recommend putting all loggers at the top of your files to register the padding properly.

goatsy supports 256 ANSI and Truecolor.

Windows CMD won’t display those colors due to it's 16 ANSI system.

```go
import (
    "time"
    "github.com/calemy/goatsy"
)

var logger = goatsy.New(&goatsy.Options{
    Truecolor: true,                       // enable Truecolor (default is 256 ANSI)
    Name: "Authentication",       // optional tag that acts as a prefix, gets padded automatically.
    TimeFormat: time.DateTime,  // optional - you can use any date format from the time package. Careful, they do not get padded.
})

var logger2 = goatsy.New(&goatsy.Options{
    Truecolor: false,
    Name: "API",
    TimeFormat: time.DateTime,
})
```

---

## Change color

You can change color with our palette of 11 colors:

Red · Orange · Yellow · Green · Aqua · Blue · Blurple · Purple · Magenta · Pink · Rose

```go
logger.Blue().
Send("you can do it like this")

logger.Blue("or like this").
Send()

logger.Blue("you even can change the color ").
Red("in the same message, ").
Green("by chaining colors together!").
Send()
```

---

## Prefixes / Names

You can even rename loggers easily after creation.

```go
logger.Rename("[Debug]").
Send("This text is a debug")
// => [Debug] | This text is a debug
logger.Rename("[Info]")..
Send("Not anymore!")
// => [Info] | Not anymore
```

Multiple loggers are aligned by longest name automatically.

```go
logger.Rename("Tim").
Send("Hey Nanoo!") 


logger2.Rename("Nanoo").
Send("Welcome back, Tim!") 


logger3.Rename("Mina").
Send("He's back already?")
// => Tim   | Hey Nanoo!
// => Nanoo | Welcome back, Tim!
// => Mina  | He's back already?
```

---

## Traces

Need to know where a log comes from? Add a stack trace!

```go
logger.Trace("Where am I?")
// => Where am I? → /home/nanoo/goatsy/main.go:10:12
```

---

## Log into files

*(Work in progress, coming back soon)*

---

## Example

Here’s a simple blue logger with the `[Info]` tag:

```go
import (
    "time"
    "github.com/calemy/goatsy"
)

var logger = goatsy.New(&goatsy.Options{
    Truecolor: true,
    Name: "[Info]",
    TimeFormat: time.DateTime
}).Blue()

logger.Send("This is supposed to serve valuable information.")
```

---

## Tips

Loggers keep the same name, timestamp, and color until you change them.

```go
logger.Blue().
Rename("Nanoo").
Send("I like the color blue")
// => Nanoo | 2025-08-18 19:44:19 | I like the color blue

logger.Send("This stays blue")
// => Nanoo | 2025-08-18 19:44:19 | This stays blue

logger.Red().
Send("until you change the color")
// => Nanoo | 2025-08-18 19:44:19 | until you change the color
```

If you find this package useful, feel free to leave a star on the repository!

Contributions are always welcomed via Issues and PRs.