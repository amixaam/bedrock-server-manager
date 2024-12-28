package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptString asks for user input with a default value
func PromptString(prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

// PromptInt asks for numeric user input with a default value
func PromptInt(prompt string, defaultValue int) int {
	input := PromptString(prompt, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return value
}

// PromptBool asks for boolean user input with a default value
func PromptBool(prompt string, defaultValue bool) bool {
	defaultStr := "no"
	if defaultValue {
		defaultStr = "yes"
	}
	input := strings.ToLower(PromptString(prompt, defaultStr))
	return input == "yes" || input == "y" || (input == "" && defaultValue)
}