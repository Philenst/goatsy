package goatsy

import (
	"fmt"
	"strconv"
)

type Logger struct {
	Truecolor bool
	messages  []Message
	color     string
}

type Message struct {
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

func (logger *Logger) Color(color string, fallback int, message ...string) *Logger {
	if logger.Truecolor {
		logger.color = color
	} else {
		logger.color = fmt.Sprintf("\x1b[38;5;%dm", fallback) //TODO: background support?
	}

	if len(message) > 0 {
		logger.messages = append(logger.messages, Message{
			Color: logger.color,
			Input: message[0],
		})
	}
	return logger
}

func (logger *Logger) Red(message ...string) *Logger {
	return logger.Color("#ff0000", 9, message...)
}

func (logger *Logger) Orange(message ...string) *Logger {
	return logger.Color("#ff8700", 208, message...)
}

func (logger *Logger) Yellow(message ...string) *Logger {
	return logger.Color("#ffff00", 226, message...)
}

func (logger *Logger) Green(message ...string) *Logger {
	return logger.Color("#00ff00", 10, message...)
}

func (logger *Logger) Aqua(message ...string) *Logger {
	return logger.Color("#00ffff", 14, message...)
}

func (logger *Logger) Blue(message ...string) *Logger {
	return logger.Color("#5f87ff", 69, message...)
}

func (logger *Logger) Blurple(message ...string) *Logger {
	return logger.Color("#5f5fff ", 63, message...)
}

func (logger *Logger) Purple(message ...string) *Logger {
	return logger.Color("#af5fff", 135, message...)
}

func (logger *Logger) Magenta(message ...string) *Logger {
	return logger.Color("#ff00ff", 13, message...)
}

func (logger *Logger) Pink(message ...string) *Logger {
	return logger.Color("#ff5faf", 205, message...)
}

func (logger *Logger) Rose(message ...string) *Logger {
	return logger.Color("#ff0087", 198, message...)
}

func (logger *Logger) Reset() *Logger {
	fmt.Print("\x1b[0m")
	return logger
}

func (logger *Logger) Send(message ...string) *Logger {
	var output string

	if len(message) > 0 {
		logger.messages = append(logger.messages, Message{
			Color: logger.color,
			Input: message[0],
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
	logger.messages = make([]Message, 0)

	output += "\x1b[0m" // Reset color at the end
	fmt.Println(output)

	return logger
}