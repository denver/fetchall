package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	csvFilename := flag.String("csv", "default.csv", "a csv file in the format of one 'http request' per line")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	links := parseLines(lines)
	fmt.Println(links)
}

func parseLines(lines [][]string) []link {
	ret := make([]link, len(lines))
	for i, line := range lines {
		ret[i] = link{
			l: line[0],
		}
	}
	return ret
}

type link struct {
	l string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
