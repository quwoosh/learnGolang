package main

import (
	"fmt"
	"io"
	"net/url"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(uri string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(uri)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	f, err := os.Create(url.QueryEscape(uri))
	if err != nil {
		ch <- err.Error()
		return
	}
	nbytes, err := io.Copy(f, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", uri, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, uri)
}