package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SavedData struct {
	value     string
	expirary  *int64
	savedTime time.Time
}

func CommandsSwitch(commands []string, savedDataMap map[string]SavedData) string {
	switch strings.ToUpper(commands[0]) {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		return fmt.Sprintf("$%v\r\n%v\r\n", len(commands[1]), commands[1])
	case "SET":
		return setCommand(commands, savedDataMap)
	case "GET":
		getCommand(commands, savedDataMap)
	}
	return "$-1\r\n"
}

func setCommand(commands []string, savedDataMap map[string]SavedData) string {
	if len(commands) < 4 {
		savedDataMap[commands[1]] = SavedData{
			value: commands[2],
		}
		return "+OK\r\n"
	}

	if strings.ToUpper(commands[3]) == "PX" {
		expiraryValue, err := strconv.ParseInt(commands[4], 10, 64)
		if err != nil {
			fmt.Println("Improper Expirary Value")
			return "$-1\r\n"
		}
		savedDataMap[commands[1]] = SavedData{
			value:     commands[2],
			expirary:  &expiraryValue,
			savedTime: time.Now(),
		}
	}
	return "+OK\r\n"
}

func getCommand(commands []string, savedDataMap map[string]SavedData) string {
	savedDataForKey := savedDataMap[commands[1]]
	if savedDataForKey.value == "" {
		return "$-1\r\n"
	}
	if savedDataForKey.expirary == nil {
		return fmt.Sprintf("$%v\r\n%v\r\n", len(savedDataForKey.value), savedDataForKey.value)
	}
	if *savedDataForKey.expirary < time.Since(savedDataForKey.savedTime).Milliseconds() {
		return fmt.Sprintf("$%v\r\n%v\r\n", len(savedDataForKey.value), savedDataForKey.value)
	}
	return "$-1\r\n"
}
