package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type result struct {
	response []byte
}

func getOne(i int, ch chan<- result) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read: %s\n", err)
		os.Exit(-1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "skipping %d: got %d\n", i, resp.StatusCode)
		ch <- result{nil}
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read from resp.Body: %s", err)
		os.Exit(-1)
	}

	ch <- result{body}
}

func getTotalNumComics() (int, error) {
	url := "https://xkcd.com/info.0.json"
	resp, err := http.Get(url)

	if err != nil {
		return 0, fmt.Errorf("can't get %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status: %d , Not found", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("can't read from resp.Body: %s", err)
	}
	var comic struct {
		Num int `json:"num"`
	}
	err = json.Unmarshal(body, &comic)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshall json: %s", err)
	}

	return comic.Num - 1, nil //-1 because comic number 404 returns status 404
}

func downloadComics() {

	var (
		output io.WriteCloser = os.Stdout
		err    error
		cnt    int
		fails  int
		data   []byte
	)

	results := make(chan result)

	if len(os.Args) > 1 {

		output, err = os.Create(os.Args[1])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		defer output.Close()
	}

	//output is in form of json array
	//add brackets before and after

	fmt.Fprintf(output, "[")
	defer fmt.Fprintf(output, "]")

	//get total number of comics from current comic
	totalNumComics, err := getTotalNumComics()

	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	for i := 1; i < totalNumComics; i++ {
		go getOne(i, results)
	}

	//stop if we get two 404s in a row
	for i := 1; i < totalNumComics; i++ {

		r := <-results

		data = r.response

		if data == nil {
			fails++
			continue
		}

		if cnt > 0 {
			fmt.Fprintf(output, ",")
		}

		_, err = io.Copy(output, bytes.NewBuffer(data))

		if err != nil {
			fmt.Fprintf(os.Stderr, "stopped: %s\n", err)
			os.Exit(1)
		}
		fails = 0
		cnt++
	}

	fmt.Fprintf(os.Stderr, "read %d comics\n", cnt)
}

func main() {
	getTotalNumComics()
	downloadComics()
}
