package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
	lane "gopkg.in/oleiade/lane.v1"
)

func main() {
	filename := flag.String("html", "", "a html file")
	flag.Parse()

	links, err := ParseAllLinks(*filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(links)
}

// ParseAllLinks parses all the links in a HTML file.
func ParseAllLinks(filename string) ([]Link, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	n, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	return parseHTMLNode(n), nil
}

func parseHTMLNode(node *html.Node) []Link {
	links := []Link{}

	s := lane.NewStack()
	s.Push(node)

	for !s.Empty() {
		t := s.Pop().(*html.Node)

		if t.Type == html.ElementNode && t.Data == "a" {
			links = append([]Link{parseLink(t)}, links...)
		}

		for c := t.FirstChild; c != nil; c = c.NextSibling {
			s.Push(c)
		}
	}

	return links
}

func parseLink(node *html.Node) Link {
	return Link{
		Href: parseLinkHref(node),
		Text: parseLinkText(node),
	}
}

func parseLinkHref(node *html.Node) string {
	for _, a := range node.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

func parseLinkText(node *html.Node) string {
	texts := []string{}

	s := lane.NewStack()
	s.Push(node)

	for !s.Empty() {
		t := s.Pop().(*html.Node)

		if t.Type == html.TextNode {
			text := strings.TrimSpace(t.Data)
			texts = append([]string{text}, texts...)
		}

		if t.Type == html.ElementNode {
			for c := t.FirstChild; c != nil; c = c.NextSibling {
				s.Push(c)
			}
		}
	}

	return strings.TrimSpace(strings.Join(texts, " "))
}

// Link represents a <link> tag
type Link struct {
	Href string
	Text string
}
