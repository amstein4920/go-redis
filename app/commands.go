package main

import (
	"fmt"
	"strings"
)

type SavedData struct {
	value string
	// expirary  *int
	// savedTime *time.Time
}

func CommandsSwitch(commands []string, savedDataMap map[string]SavedData) string {
	switch strings.ToUpper(commands[0]) {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		return fmt.Sprintf("$%v\r\n%v\r\n", len(commands[1]), commands[1])
	case "SET":
		if len(commands) < 4 {
			savedDataMap[commands[1]] = SavedData{
				value: commands[2],
			}
			break
		}
		return "+OK\r\n"
	case "GET":
		savedDataForKey := savedDataMap[commands[1]]
		if savedDataForKey.value != "" {
			return fmt.Sprintf("$%v\r\n%v\r\n", len(savedDataForKey.value), savedDataForKey.value)
		}
	}
	return "$-1\r\n"
}
