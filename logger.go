package goatsy

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var padding int = 0
var names = map[string]int{}

type Options struct {
	Truecolor  bool
	Name       string
	TimeFormat string
}

type Logger struct {
	truecolor  bool
	name       string
	timeFormat string
	messages   []message
	color      string
}

type message struct {
	Color string
	Input string
}

func New(options *Options) *Logger {
	if options.Name != "" {
		if len(options.Name) > padding {
			padding = len(options.Name)
		}
	}
	names[options.Name] = len(options.Name)
	return &Logger{
		truecolor:  options.Truecolor,
		name:       options.Name,
		timeFormat: options.TimeFormat,
	}
}

func convertHex(hex string) string {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)

	prefix := 38
	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", prefix, r, g, b)
}

func (l *Logger) Color(color string, fallback int, input ...string) *Logger {
	if l.truecolor {
		l.color = color
	} else {
		l.color = fmt.Sprintf("\x1b[38;5;%dm", fallback) //TODO: background support?
	}

	if len(input) > 0 {
		l.messages = append(l.messages, message{
			Color: l.color,
			Input: input[0],
		})
	}
	return l
}

func (l *Logger) Red(input ...string) *Logger {
	return l.Color("#ff0000", 9, input...)
}

func (l *Logger) Orange(input ...string) *Logger {
	return l.Color("#ff8700", 208, input...)
}

func (l *Logger) Yellow(input ...string) *Logger {
	return l.Color("#ffff00", 226, input...)
}

func (l *Logger) Green(input ...string) *Logger {
	return l.Color("#00ff00", 10, input...)
}

func (l *Logger) Aqua(input ...string) *Logger {
	return l.Color("#00ffff", 14, input...)
}

func (l *Logger) Blue(input ...string) *Logger {
	return l.Color("#5f87ff", 69, input...)
}

func (l *Logger) Blurple(input ...string) *Logger {
	return l.Color("#5f5fff", 63, input...)
}

func (l *Logger) Purple(input ...string) *Logger {
	return l.Color("#af5fff", 135, input...)
}

func (l *Logger) Magenta(input ...string) *Logger {
	return l.Color("#ff00ff", 13, input...)
}

func (l *Logger) Pink(input ...string) *Logger {
	return l.Color("#ff5faf", 205, input...)
}

func (l *Logger) Rose(input ...string) *Logger {
	return l.Color("#ff0087", 198, input...)
}

func (l *Logger) Reset() *Logger {
	fmt.Print("\x1b[0m")
	return l
}

func scan() {
	h := 0
	for _, v := range names {
		if v > h {
			h = v
		}
	}
	padding = h
}

func (l *Logger) Destroy() {
	delete(names, l.name)
	if len(l.name) == padding {
		scan()
	}
}

func (l *Logger) send(traced bool, input ...string) *Logger {
	var output string

	if len(input) > 0 {
		l.messages = append(l.messages, message{
			Color: l.color,
			Input: input[0],
		})
	}

	if l.name != "" {
		if l.truecolor {
			output += convertHex(l.color) + fmt.Sprintf("%-*s | ", padding, l.name)
		} else {
			output += l.color + fmt.Sprintf("%-*s | ", padding, l.name)
		}
	}

	if l.timeFormat != "" {
		if l.truecolor {
			output += convertHex(l.color) + fmt.Sprintf("%s | ", time.Now().Format(l.timeFormat))
		} else {
			output += l.color + fmt.Sprintf("%s | ", time.Now().Format(l.timeFormat))
		}
	}

	for _, msg := range l.messages {
		if l.truecolor {
			output += convertHex(msg.Color) + msg.Input
		} else {
			output += msg.Color + msg.Input
		}
	}

	l.messages = nil
	l.messages = make([]message, 0)

	if traced {
		_, file, line, ok := runtime.Caller(2)

		if !ok {
			file, line = "Unknown", 0
		}

		if l.truecolor {
			output += convertHex(l.color) + fmt.Sprintf(" → %s:%d", file, line)
		} else {
			output += l.color + fmt.Sprintf(" → %s:%d", file, line)
		}

	}

	output += "\x1b[0m" // Reset color at the end
	fmt.Println(output)
	return l
}

func (l *Logger) Send(input ...string) *Logger {
	return l.send(false, input...)
}

func (l *Logger) Trace(input ...string) *Logger {
	return l.send(true, input...)
}

func (l *Logger) Rename(name string) *Logger {
	delete(names, l.name)
	names[name] = len(name)

	if len(name) >= padding {
		l.name = name
		padding = len(name)
		return l
	}

	if len(l.name) == padding {
		l.name = name
		scan()
	}
	return l
}
