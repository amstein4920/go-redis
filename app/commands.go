package main

import (
	"fmt"
	"strings"
)

func CommandsSwitch(commands []string, savedDataMap map[string]string) string {
	switch strings.ToUpper(commands[0]) {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		return fmt.Sprintf("$%v\r\n%v\r\n", len(commands[1]), commands[1])
	case "SET":
		savedDataMap[commands[1]] = commands[2]
		return "+OK\r\n"
	case "GET":
		valueForKey := savedDataMap[commands[1]]
		if valueForKey != "" {
			return fmt.Sprintf("$%v\r\n%v\r\n", len(valueForKey), valueForKey)
		}
	}
	return "$-1\r\n"
}
