package goatsy

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var padding int = 0
var names = map[string]int{}
var mu sync.Mutex

type Options struct {
	Truecolor  bool
	Name       string
	TimeFormat string
	LogPath    string
}

type Logger struct {
	mu         sync.Mutex
	truecolor  bool
	name       string
	timeFormat string
	messages   []message
	color      string
	file       *os.File
}

type message struct {
	Color string
	Input string
}

type entry struct {
	name       string
	timeFormat string
	truecolor  bool
	color      string
	messages   []message
	file       *os.File
}

func New(options *Options) *Logger {
	mu.Lock()
	defer mu.Unlock()

	if options.Name != "" && len(options.Name) > padding {
		padding = len(options.Name)
	}
	names[options.Name] = len(options.Name)

	var file *os.File
	if options.LogPath != "" {
		_ = os.MkdirAll(filepath.Dir(options.LogPath), 0755)
		f, err := os.OpenFile(options.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			file = f
		}
	}

	return &Logger{
		truecolor:  options.Truecolor,
		name:       options.Name,
		timeFormat: options.TimeFormat,
		file:       file,
	}
}

func convertHex(hex string) string {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func (l *Logger) flush(input ...string) entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(input) > 0 {
		l.messages = append(l.messages, message{
			Color: l.color,
			Input: input[0],
		})
	}

	msgs := append([]message(nil), l.messages...)
	e := entry{
		name:       l.name,
		timeFormat: l.timeFormat,
		truecolor:  l.truecolor,
		color:      l.color,
		messages:   msgs,
		file:       l.file,
	}

	l.messages = nil
	return e
}

func formatConsole(e entry) string {
	var output string

	if e.name != "" {
		mu.Lock()
		pad := padding
		mu.Unlock()

		if e.truecolor {
			output += convertHex(e.color) + fmt.Sprintf("%-*s | ", pad, e.name)
		} else {
			output += e.color + fmt.Sprintf("%-*s | ", pad, e.name)
		}
	}

	if e.timeFormat != "" {
		ts := time.Now().Format(e.timeFormat)
		if e.truecolor {
			output += convertHex(e.color) + fmt.Sprintf("%s | ", ts)
		} else {
			output += e.color + fmt.Sprintf("%s | ", ts)
		}
	}

	for _, msg := range e.messages {
		if e.truecolor {
			output += convertHex(msg.Color) + msg.Input
		} else {
			output += msg.Color + msg.Input
		}
	}

	output += "\x1b[0m"
	return output
}

func formatFile(e entry) string {
	var output string

	if e.timeFormat != "" {
		output += fmt.Sprintf("%s | ", time.Now().Format(e.timeFormat))
	}

	for _, msg := range e.messages {
		output += msg.Input
	}

	return output
}

func (l *Logger) Color(color string, fallback int, input ...string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.truecolor {
		l.color = color
	} else {
		l.color = fmt.Sprintf("\x1b[38;5;%dm", fallback)
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

func (l *Logger) Save(input ...string) *Logger {
	e := l.flush(input...)
	if e.file != nil {
		_, _ = e.file.WriteString(formatFile(e) + "\n")
	}
	return l
}

func (l *Logger) Log(input ...string) *Logger {
	e := l.flush(input...)
	fmt.Println(formatConsole(e))
	if e.file != nil {
		_, _ = e.file.WriteString(formatFile(e) + "\n")
	}
	return l
}

func (l *Logger) Reset() *Logger {
	fmt.Print("\x1b[0m")
	return l
}

func (l *Logger) Destroy() {
	l.mu.Lock()
	if l.file != nil {
		_ = l.file.Close()
		l.file = nil
	}
	l.mu.Unlock()

	mu.Lock()
	defer mu.Unlock()

	delete(names, l.name)
	if len(l.name) == padding {
		h := 0
		for _, v := range names {
			if v > h {
				h = v
			}
		}
		padding = h
	}
}

func (l *Logger) Rename(name string) *Logger {
	l.mu.Lock()
	mu.Lock()
	defer mu.Unlock()
	defer l.mu.Unlock()

	delete(names, l.name)
	names[name] = len(name)
	l.name = name

	if len(name) > padding {
		padding = len(name)
	}

	return l
}