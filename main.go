package main

import (
	"fmt"
	"bufio"
	"os"
	"os/exec"
	"strings"
)

var cliName string = "RSSFeedReader"

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func printUnknown(text string) {
	fmt.Println(text, ": command not found")
}

func displayHelp() {
	fmt.Println("help")
}

func clearScreen() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func handleInvalidCmd(text string) {
	defer printUnknown(text)
}

func handleCmd(text string) {
	handleInvalidCmd(text)
}

func cleanInput(text string) string {
	output := strings.TrimSpace(text)
	output = strings.ToLower(output)
	return output
}

func formatInput(text string) []string {
	var querry []string

	size := len(text)

	left := 0

	for i := range size {
		if text[i] == ' ' {
			querry = append(querry, text[left:i])
			left = i + 1
		}
	}
	querry = append(querry, text[left:size])

	return querry
}

func main() {
	commands := map[string]interface{} {
		".help": displayHelp,
		".clear": clearScreen,
	}

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()

	for reader.Scan() {

		input := cleanInput(reader.Text())
		var querry []string = formatInput(input)

		var cmd string = querry[0]

		if command, exists := commands[cmd]; exists {
			command.(func()) ()
		} else if strings.EqualFold(".exit", cmd) {
			return
		} else {
			handleCmd(cmd)
		}
		printPrompt()
	}
	fmt.Println()
}
