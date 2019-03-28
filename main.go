package main

/*
A simple program that takes in valid linux files or paths and checks whether they include critical or sensitive files.
*/

// TODO: handling of filepaths
// TODO: Stdin inputs for files
// TODO: extension matching
// TODO: cleanup of rules file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type rule struct {
	Application string
	Description string
	Name        string
	Pattern     string
	Regex       string
	Part        string
}

var wg sync.WaitGroup

func openFiles() []string { // Will be replaced with reading from stdin.
	var files []string
	file, err := os.Open("/home/calmus/Tech/golang/src/github.com/Calmusroot/Grepper/files.txt") // open	fmt.Println("asd")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewScanner(file)
	defer file.Close()
	for reader.Scan() {
		files = append(files, reader.Text())
	}
	return files
}

func loadRules() []rule {
	file, err := os.Open("/home/calmus/Tech/golang/src/github.com/Calmusroot/Grepper/rules.csv")
	if err != nil {
		panic(err)
	}
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

func checkFiles(paths []string, r rule) {
	defer wg.Done()
	for _, path := range paths {
		file := filepath.Base(path)
		m, err := regexp.MatchString(r.Pattern, file)
		if m {
			fmt.Println("[*]", r.Name, "File:", path, err)
		}
	}
}

func main() {
	paths := openFiles()
	rules := loadRules()
	// messages := make(chan string, len(rules))
	for _, r := range rules {
		// fmt.Println(r)
		wg.Add(1)
		go checkFiles(paths, r)
	}
	wg.Wait()
	// close(messages)
}
