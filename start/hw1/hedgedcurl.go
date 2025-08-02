package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	timeout  int
	showHelp bool
)

func init() {
	flag.IntVar(&timeout, "timeout", 15, "максимальное время ожидания ответа")
	flag.IntVar(&timeout, "t", 15, "сокращение для --timeout")
	flag.BoolVar(&showHelp, "help", false, "показать справку")
	flag.BoolVar(&showHelp, "h", false, "сокращение для --help")
}

func printHelp() {
	fmt.Printf("Использование: %s [OPTIONS] URL1 URL2...\n\n", os.Args[0])
	fmt.Println("Опции:")
	flag.PrintDefaults()
	fmt.Println("\nПримеры:")
	fmt.Printf("  %s -t 3s https://google.com https://ya.ru\n", os.Args[0])
	fmt.Printf("  %s --timeout=1s https://example.com\n", os.Args[0])
}

func fetchURL(ctx context.Context, url string, result chan<- *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	select {
	case result <- resp:
	case <-ctx.Done():
		resp.Body.Close()
	}
}

func printResponse(resp *http.Response) {
	fmt.Printf("HTTP %s\n", resp.Status)

	for k, v := range resp.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
	} else {
		fmt.Println(string(body))
	}
	resp.Body.Close()
}

func main() {
	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	urls := flag.Args()
	if len(urls) == 0 {
		printHelp()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := make(chan *http.Response, 1)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go fetchURL(ctx, url, result, &wg)
	}

	select {
	case resp := <-result:
		printResponse(resp)
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Print(228)
		os.Exit(228)
	}

	cancel()
	wg.Wait()
}
