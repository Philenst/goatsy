package goatsy

import (
	"fmt"
	"strconv"
)

type Logger struct {
	Truecolor bool
	messages  []message
	color     string
}

type message struct {
	Color string
	Input string
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

func (logger *Logger) Color(color string, fallback int, input ...string) *Logger {
	if logger.Truecolor {
		logger.color = color
	} else {
		logger.color = fmt.Sprintf("\x1b[38;5;%dm", fallback) //TODO: background support?
	}

	if len(input) > 0 {
		logger.messages = append(logger.messages, message{
			Color: logger.color,
			Input: input[0],
		})
	}
	return logger
}

func (logger *Logger) Red(input ...string) *Logger {
	return logger.Color("#ff0000", 9, input...)
}

func (logger *Logger) Orange(input ...string) *Logger {
	return logger.Color("#ff8700", 208, input...)
}

func (logger *Logger) Yellow(input ...string) *Logger {
	return logger.Color("#ffff00", 226, input...)
}

func (logger *Logger) Green(input ...string) *Logger {
	return logger.Color("#00ff00", 10, input...)
}

func (logger *Logger) Aqua(input ...string) *Logger {
	return logger.Color("#00ffff", 14, input...)
}

func (logger *Logger) Blue(input ...string) *Logger {
	return logger.Color("#5f87ff", 69, input...)
}

func (logger *Logger) Blurple(input ...string) *Logger {
	return logger.Color("#5f5fff", 63, input...)
}

func (logger *Logger) Purple(input ...string) *Logger {
	return logger.Color("#af5fff", 135, input...)
}

func (logger *Logger) Magenta(input ...string) *Logger {
	return logger.Color("#ff00ff", 13, input...)
}

func (logger *Logger) Pink(input ...string) *Logger {
	return logger.Color("#ff5faf", 205, input...)
}

func (logger *Logger) Rose(input ...string) *Logger {
	return logger.Color("#ff0087", 198, input...)
}

func (logger *Logger) Reset() *Logger {
	fmt.Print("\x1b[0m")
	return logger
}

func (logger *Logger) Send(input ...string) *Logger {
	var output string

	if len(input) > 0 {
		logger.messages = append(logger.messages, message{
			Color: logger.color,
			Input: input[0],
		})
	}

	for _, msg := range logger.messages {
		if logger.Truecolor {
			output += convertHex(msg.Color) + msg.Input
		} else {
			output += msg.Color + msg.Input
		}
	}

	logger.messages = nil
	logger.messages = make([]message, 0)

	output += "\x1b[0m" // Reset color at the end
	fmt.Println(output)

	return logger
}