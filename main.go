package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strings"
)

type Feed struct {
	Site string `json:"site"`
	Link string `json:"link"`
	Type string `json:"type"`
}

func readRSS(f string) []string {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urls := getURLs(scanner)

	return urls
}

func getURLs(scanner *bufio.Scanner) []string {
	var urls []string

	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), " ", 3)[2:]
		url := strings.Join(line, "")
		urls = append(urls, url)
	}
	// https://go.dev/blog/errors-are-values
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return urls
}

func stripHost(host string) string {
	prefixes := []string{"www.", "feeds.", "rss.", "rssfeeds."}
	suffixes := []string{".com", ".co.uk", ".com.au", ".net.au", ".co.nz", ".yahoo", ".in", ".feedsportal", ".ca", ".net", ".org", ".gov"}
	for _, prefix := range prefixes {
		host = strings.TrimPrefix(host, prefix)
	}
	for _, suffix := range suffixes {
		host = strings.TrimSuffix(host, suffix)
	}

	return host
}

func getFeeds(urls []string) []Feed {
	feeds := []Feed{}
	for _, s := range urls {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		host := stripHost(u.Host)
		feed := Feed{
			Site: host,
			Link: s,
			Type: "rss",
		}
		feeds = append(feeds, feed)
	}

	return feeds
}

func writeJson(fileW string, feeds []Feed) {
	file, _ := json.MarshalIndent(feeds, "", "  ")

	err := os.WriteFile(fileW, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success")
}

func main() {
	const srcRSS = "data/rss.txt"
	const dstRSS = "data/rss.json"

	urls := readRSS(srcRSS)
	feeds := getFeeds(urls)
	writeJson(dstRSS, feeds)
}
