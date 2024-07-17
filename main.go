package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

func handleInvalidCmd(err string) {
	defer printUnknown(querry[0], err)
}

func handleCmd() {
	handleInvalidCmd("command not found")
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
		handleInvalidCmd("invalid url")
		return false
	}
	return true
}

func addLink() {
	size := len(querry)
	if size != 2 && size != 3 {
		handleInvalidCmd("invalid number of arguments")
		return
	}
	var s string = querry[1]
	u, _ := url.Parse(s)
	var category string

	if size == 2 {
		category = "uncategorized"
	} else {
		category = querry[2]
	}

	if !validateUrl(s) {
		handleInvalidCmd("invalid url")
		return
	}
	link_map[category] = append(link_map[category], feed{u.Host, s})
}

func removeLink() {
	size := len(querry)
	if size != 2 {
		handleInvalidCmd("invalid number of arguments")
		return
	}
	var feedToRemove string = querry[1]
	for key, arr := range link_map {
		newArr := []feed{}
		for _, fd := range arr {
			if !strings.EqualFold(feedToRemove, fd.Name) {
				newArr = append(newArr, fd)
			}
		}
		if len(newArr) == 0 {
			delete(link_map, key)
		} else {
			link_map[key] = newArr
		}
	}
}

func changeCategory() {
	size := len(querry)
	if size != 3 {
		handleInvalidCmd("invalid number of arguments")
		return
	}

	name := querry[1]
	category := querry[2]

	var link string
	found := false

	for key, arr := range link_map {
		for i, fd := range arr {
			if strings.EqualFold(name, fd.Name) {
				link = fd.Link
				link_map[key] = append(arr[:i], arr[i+1:]...)
				if len(link_map[key]) == 0 {
					delete(link_map, key)
				}
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		handleInvalidCmd("feed not found")
		return
	}

	querry = []string{".add", link, category}
	addLink()

}

func showLinks() {
	size := len(querry)
	if size != 1 && size != 2 {
		handleInvalidCmd("invalid number of arguments")
		return
	}
	fmt.Println()
	if size == 2 {
		category := querry[1]
		fmt.Printf("%s:\n", category)
		for idx, itm := range link_map[category] {
			ch := ";"
			if idx == len(link_map[category]) - 1 {
				ch = "."
			} 
			fmt.Printf("    %d. %s %s\n", idx + 1, itm.Name, ch)
		}

		fmt.Println()
	} else {
		for category, arr := range link_map {
			fmt.Printf("%s:\n", category)
			for idx, itm := range arr {
				ch := ";"
				if idx == len(arr) - 1 {
					ch = "."
				} 
				fmt.Printf("    %d. %s %s\n", idx + 1, itm.Name, ch)
			}
			fmt.Println()
		}
	}
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

func openFeed() {
	size := len(querry)
	if size != 2 {
		handleInvalidCmd("invalid number of arguments")
		return
	}

	name := querry[1]
	var link string
	found := false

	for _, arr := range link_map {
		for i := range arr {
			if strings.EqualFold(name, arr[i].Name) {
				link = arr[i].Link
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		handleInvalidCmd("feed not found")
		return
	}

	resp, httpErr := http.Get(link)
	if httpErr != nil {
		handleInvalidCmd("error getting the data")
		return
	}

	body, readingErr := ioutil.ReadAll(resp.Body)
	if readingErr != nil {
		handleInvalidCmd("error reading the data")
		return
	}

	fmt.Println(string(body))
}

type feed struct {
	Name string
	Link string
}

var link_map map[string][]feed

var querry []string

var commands = map[string]interface{} {
	".help": displayHelp,
	".clear": clearScreen,
	".add": addLink,
	".remove": removeLink,
	".category": changeCategory,
	".show": showLinks,
	".open": openFeed,
}

func main() {
	loadJsonFile();

	//fmt.Println(link_map)

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
			handleCmd()
		}
		printPrompt()
	}
	fmt.Println()

	saveLinksToFile()
}
