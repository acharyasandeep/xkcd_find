# XKCD Comics Download and Search

#### A concurrent xkcd comics downloader and searcher in golang.

## Setup

- Open terminal and `cd` to `downloader` directory 
- Run `go run xkcd_download.go xkcd.json` 
- `cd` to `../searcher` directory and search for keywords in comics using `go run searcher.go ../downloader/xkcd.json keyword_to_search`