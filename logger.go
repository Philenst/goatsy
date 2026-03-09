package goatsy

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	mu          sync.Mutex
	truecolor   bool
	name        string
	timeFormat  string
	messages    []message
	color       string
	logPath     string
	file        *os.File
	currentDate string
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
}

func New(options *Options) *Logger {
	mu.Lock()
	defer mu.Unlock()

	if options.Name != "" && len(options.Name) > padding {
		padding = len(options.Name)
	}
	names[options.Name] = len(options.Name)

	return &Logger{
		truecolor:  options.Truecolor,
		name:       options.Name,
		timeFormat: options.TimeFormat,
		logPath:    options.LogPath,
	}
}

func datedLogPath(basePath, date string) string {
	ext := filepath.Ext(basePath)
	name := strings.TrimSuffix(filepath.Base(basePath), ext)
	dir := filepath.Dir(basePath)

	if ext == "" {
		ext = ".log"
	}

	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, date, ext))
}

func (l *Logger) ensureLogFile(now time.Time) *os.File {
	if l.logPath == "" {
		return nil
	}

	date := now.Format("2006-01-02")
	if l.file != nil && l.currentDate == date {
		return l.file
	}

	if l.file != nil {
		_ = l.file.Close()
		l.file = nil
	}

	fullPath := datedLogPath(l.logPath, date)
	_ = os.MkdirAll(filepath.Dir(fullPath), 0755)

	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}

	l.file = f
	l.currentDate = date
	return l.file
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
	}

	l.messages = nil
	return e
}

func formatConsole(e entry, now time.Time) string {
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
		ts := now.Format(e.timeFormat)
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

func formatFile(e entry, now time.Time) string {
	var output string

	if e.timeFormat != "" {
		output += fmt.Sprintf("%s | ", now.Format(e.timeFormat))
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
	now := time.Now()
	e := l.flush(input...)

	l.mu.Lock()
	file := l.ensureLogFile(now)
	l.mu.Unlock()

	if file != nil {
		_, _ = file.WriteString(formatFile(e, now) + "\n")
	}

	return l
}

func (l *Logger) Log(input ...string) *Logger {
	now := time.Now()
	e := l.flush(input...)

	fmt.Println(formatConsole(e, now))

	l.mu.Lock()
	file := l.ensureLogFile(now)
	l.mu.Unlock()

	if file != nil {
		_, _ = file.WriteString(formatFile(e, now) + "\n")
	}

	return l
}

func (l *Logger) Send(input ...string) *Logger {
	now := time.Now()
	e := l.flush(input...)
	fmt.Println(formatConsole(e, now))
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