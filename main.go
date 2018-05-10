package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
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
	start := time.Now()
	ch := make(chan string)
	for _, url := range links {
		// fmt.Printf("%d: %s\n", i+1, p.l)
		go fetch(url.l, ch)
	}
	for range links {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v ", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

func parseLines(lines [][]string) []link {
	ret := make([]link, len(lines))
	for i, line := range lines {
		ret[i] = link{
			l: strings.TrimSpace(line[0]),
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
