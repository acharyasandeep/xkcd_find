package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type xkcdComic struct {
	Num        int    `json:"num"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Img        string `json:"img"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no file given")
		os.Exit(1)
	}

	fn := os.Args[1]

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "no search term")
		os.Exit(-1)
	}

	var (
		items []xkcdComic
		terms []string
		input io.ReadCloser
		cnt   int
		err   error
	)
	if input, err = os.Open("../downloader/" + fn); err != nil {
		fmt.Fprintln(os.Stderr, "Can't open file:", err)
		os.Exit(-1)
	}

	defer input.Close()

	//decode file

	if err = json.NewDecoder(input).Decode(&items); err != nil {
		fmt.Fprintln(os.Stderr, "Can't decode json:", err)
		os.Exit(-1)
	}

	fmt.Fprintf(os.Stderr, "read %d comics\n", len(items))

	//get search terms
	for _, t := range os.Args[2:] {
		terms = append(terms, strings.ToLower(t))
	}

	//search
outer:
	for _, item := range items {
		title := strings.ToLower(item.Title)
		transcript := strings.ToLower(item.Transcript)

		for _, term := range terms {
			if !strings.Contains(title, term) && !strings.Contains(transcript, term) {
				continue outer
			}
		}
		fmt.Printf("https://xkcd.com/%d/ %s/%s/%s %q\n",
			item.Num, item.Month, item.Day, item.Year, item.Title)
		cnt++
	}

	fmt.Fprintf(os.Stderr, "found %d comics\n", cnt)

}
