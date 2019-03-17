package main

/*\
A simple program that takes in valid linux files or paths and checks whether they include critical or sensitive files.
*/

// TODO: handling of filepaths
// TODO: goroutines for regex
// TODO: Stdin inputs for files
// TODO: extension matching
// TODO: cleanup of rules file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type rule struct {
	Application string
	Description string
	Name        string
	Pattern     string
	Regex       string
	Part        string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func openFiles() []string { // will be replaced with reading form stdin
	var files []string // declare a slice to hold all files

	file, err := os.Open("/home/calmus/Tech/golang/src/github.com/Calmusroot/Grepper/files.txt") // open	fmt.Println("asd")
	check(err)
	reader := bufio.NewScanner(file)
	defer file.Close()
	for reader.Scan() {
		files = append(files, reader.Text())
	}
	return files
}
func checkRegex(rgx *regexp.Regexp, file string) bool {
	if rgx.MatchString(file) {
		return true
	}
	return false
}
func checkFiles(files []string, r rule, messages chan string) {
	if r.Application == "regex" {
		rgx, _ := regexp.Compile(r.Pattern)
		for _, file := range files {
			if checkRegex(rgx, file) {
				// messages <- "[*] Found:", r.Name, "File:", file
				messages <- "[*] Found:"
			}
		}
	}
}

func loadRules() []rule {
	file, err := os.Open("/home/calmus/Tech/golang/src/github.com/Calmusroot/Grepper/rules.csv")
	check(err)
	reader := csv.NewReader(file)
	defer file.Close()
	var rules []rule
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		rules = append(rules, rule{
			Part:        strings.TrimSpace(line[0]),
			Application: strings.TrimSpace(line[1]),
			Pattern:     strings.TrimSpace(line[2]),
			Name:        strings.TrimSpace(line[3]),
			Description: strings.TrimSpace(line[4]),
		})
	}
	return rules
}

func main() {
	files := openFiles()
	rules := loadRules()
	messages := make(chan string)
	for _, r := range rules {
		go checkFiles(files, r, messages)
	}
	msg := <-messages
	fmt.Println(msg)
}
