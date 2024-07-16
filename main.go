package main

import (
	"fmt"
	"bufio"
	"os"
	"os/exec"
	"strings"
	"net/url"
	"encoding/json"
	"io/ioutil"
)

var cliName string = "RSSFeedReader"
var linkSaveFile string = "./links.json"

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func printUnknown(cmd string, err string) {
	fmt.Println(cmd, ":", err)
}

func displayHelp() {
	fmt.Println("help")
}

func clearScreen() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func handleInvalidCmd(cmd string, err string) {
	defer printUnknown(cmd, err)
}

func handleCmd(cmd string) {
	handleInvalidCmd(cmd, "command not found")
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

func validateUrl(text string) bool {
	_, err := url.ParseRequestURI(text)
	if err != nil {
		handleInvalidCmd(".add", "invalid url")
		return false
	}
	return true
}

func addLink() {
	size := len(querry)
	if size != 2 && size != 3 {
		handleInvalidCmd(".add", "invalid number of arguments")
		return
	}
	var url string = querry[1]
	var category string

	if size == 2 {
		category = "uncategorized"
	} else {
		category = querry[2]
	}

	if validateUrl(url) {
		link_map[category] = append(link_map[category], url)
	}
}

func removeLink() {
	size := len(querry)
	if size != 2 {
		handleInvalidCmd(".remove", "invalid number of arguments")
		return
	}
	var url_to_be_deleted string = querry[1]
	for key, arr := range link_map {
		newArr := arr[:0]
		for _, url := range arr {
			if !strings.EqualFold(url_to_be_deleted, url) {
				newArr = append(newArr, url)
			}
		}
		link_map[key] = newArr
		if len(newArr) == 0 {
			delete(link_map, key)
		}
	}
}

func changeCategory() {
	size := len(querry)
	if size != 3 {
		handleInvalidCmd(".category", "invalid number of arguments")
		return
	}

	category := querry[2]

	newQuerry := querry[:2]
	querry = newQuerry

	removeLink()

	querry = append(querry, category)
	addLink()
}

func showLinks() {
	size := len(querry)
	if size != 1 && size != 2 {
		handleInvalidCmd(".show", "invalid number of arguments")
		return
	}
	fmt.Println()
	if size == 2 {
		category := querry[1]
		fmt.Println()
		fmt.Printf("%s:\n", category)
		for idx, itm := range link_map[category] {
			ch := ";"
			if idx == len(link_map[category]) - 1 {
				ch = "."
			} 
			fmt.Printf("    %d. %s %s\n", idx + 1, itm, ch)
		}

		fmt.Println()
	} else {
		for category, arr := range link_map {
			fmt.Println()
			fmt.Printf("%s:\n", category)
			for idx, itm := range arr {
				ch := ";"
				if idx == len(arr) - 1 {
					ch = "."
				} 
				fmt.Printf("    %d. %s %s\n", idx + 1, itm, ch)
			}
			fmt.Println()
		}
	}
	fmt.Println()
}

func loadJsonFile() {
	content, err := ioutil.ReadFile(linkSaveFile )
	if err == nil {
		json.Unmarshal(content, &link_map)
	}
}

func saveLinksToFile() {
	asJSON, err := json.MarshalIndent(link_map, "", "\t")
	if err != nil {
		fmt.Printf("Error marshalling data to JSON: %v\n", err)
		return
	}

	// Write the JSON data to the specified file
	writingErr := ioutil.WriteFile(linkSaveFile, asJSON, 0644)
	if writingErr != nil {
		fmt.Printf("Error writing to file: %v\n", writingErr)
	}
}

var link_map = make(map[string][]string)
var querry []string

var commands = map[string]interface{} {
	".help": displayHelp,
	".clear": clearScreen,
	".add": addLink,
	".remove": removeLink,
	".category": changeCategory,
	".show": showLinks,
}

func main() {
	loadJsonFile();

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()

	for reader.Scan() {

		input := cleanInput(reader.Text())
		querry = formatInput(input)

		var cmd string = querry[0]

		if command, exists := commands[cmd]; exists {
			command.(func()) ()
		} else if strings.EqualFold(".exit", cmd) {
			saveLinksToFile()
			return
		} else {
			handleCmd(cmd)
		}
		printPrompt()
	}
	fmt.Println()

	saveLinksToFile()
}
