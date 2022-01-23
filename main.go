package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urls := scanner.Text()
		dotgit := urls + "/.git/config"
		resp, err := http.Get(dotgit)
		if err != nil {
			continue
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
}
