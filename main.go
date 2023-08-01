package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func checkGitDirectory(urls string, wg *sync.WaitGroup) {
	defer wg.Done()

	dotgit := urls + "/.git/config"
	resp, err := http.Get(dotgit)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	respString := string(body)
	if strings.Contains(respString, "repositoryformatversion") {
		fmt.Printf("[\033[32mVulnerable\033[0m] %s\n", urls)
	} else {
		fmt.Printf("[\033[31mNot Vulnerable\033[0m] %s\n", urls)
	}
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 30, "concurrency count")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [filename]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	var inputReader io.Reader
	args := flag.Args()
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()
		inputReader = file
	} else {
		inputReader = os.Stdin
	}

	scanner := bufio.NewScanner(inputReader)

	var wg sync.WaitGroup

	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		go func(u string) {
			checkGitDirectory(u, &wg)
		}(url)
	}

	wg.Wait()
}
