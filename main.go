package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type URLFormatter interface {
	Format(rawurl string) (string, error)
}

type MarkdownURLFormatter struct {
}

func (f MarkdownURLFormatter) fetchTitle(rawurl string) (string, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	return doc.Find("title").First().Text(), nil
}

func (f MarkdownURLFormatter) Format(rawurl string) (string, error) {
	title, err := f.fetchTitle(rawurl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch title: %w", err)
	}
	return fmt.Sprintf("[%s](%s)", title, rawurl), nil
}

func main() {
	flagURL := flag.String("url", "", "URL for fetch")
	flagFormatter := flag.String("formatter", "md", "URL formatter")

	flag.Parse()

	if flagURL == nil || *flagURL == "" {
		log.Fatalf("url must not to be empty: -url=%#v", flagURL)
	}

	var formatter URLFormatter
	switch *flagFormatter {
	case "md":
		formatter = MarkdownURLFormatter{}
	default:
		log.Fatalf("invalid argument: -formatter=%#v", flagFormatter)
	}

	formatted, err := formatter.Format(*flagURL)
	if err != nil {
		log.Fatalf("failed to format URL: %v", err)
	}
	fmt.Println(formatted)
}
