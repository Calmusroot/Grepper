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
	File        string
}

func (r *rule) SetFile(file string) {
	r.File = file
}

var wg sync.WaitGroup
var fileQueue sync.WaitGroup

func check(e error) { //Implement a generic error function.
	if e != nil {
		panic(e)
	}
}

func openFiles() []string { // Will be replaced with reading from stdin.
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
		// fmt.Println(line)
		rules = append(rules, rule{
			Part:        strings.TrimSpace(line[0]),
			Application: strings.TrimSpace(line[1]),
			Pattern:     strings.TrimSpace(line[2]),
			Name:        strings.TrimSpace(line[3]),
			Description: strings.TrimSpace(line[4]),
			File:        "a",
		})
	}
	return rules
}

func worker(id int, jobs <-chan rule, results chan<- string) {
	for job := range jobs {
		if job.Application == "regex" {
			rgx, _ := regexp.Compile(job.Pattern)
			if rgx.MatchString(job.File) {
				s := job.File + job.Description
				fmt.Println(job)
				results <- s
			}
		}
		results <- "a"
	}
}

func main() {
	files := openFiles()
	rules := loadRules()

	jobs := make(chan rule, len(rules)*len(files))
	results := make(chan string, len(rules)*len(files))

	for j := 0; j < 100; j++ {
		go worker(j, jobs, results)
	}
	fmt.Println("Worker spwanded")
	for _, r := range rules {
		for _, file := range files {
			r.SetFile(file)
			jobs <- r
		}
	}
	fmt.Println("Jobs out")
	close(jobs)
	for j := 0; j < len(rules)*len(files); j++ {
		<-results
		// fmt.Println(strings.TrimSuffix(<-results, "a"))
	}
}
