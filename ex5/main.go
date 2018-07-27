package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jackytck/gophercises/ex4/link"
	lane "gopkg.in/oleiade/lane.v1"
)

func main() {
	domain := flag.String("domain", "https://www.golang.org", "Root domain")
	depth := flag.Int("depth", 1, "Depth of transversal")
	out := flag.String("out", "sitemap.xml", "Output filename")
	flag.Parse()

	urls := crawl(*domain, *depth)
	if err := generateXML(urls, *out); err != nil {
		panic(err)
	}
}

func generateXML(urls urlset, out string) error {
	f, err := os.Create(out)
	defer f.Close()
	if err != nil {
		return err
	}
	urls.XMLns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	f.Write([]byte(xml.Header))
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err := enc.Encode(&urls); err != nil {
		return err
	}
	return nil
}

type urlset struct {
	XMLName xml.Name `xml:"urlset"`
	XMLns   string   `xml:"xmlns,attr"`
	URLs    []url
}

type url struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

type node struct {
	Link  link.Link
	Depth int
}

func crawl(domain string, depth int) urlset {
	root := node{
		Link: link.Link{
			Href: domain,
		},
		Depth: -1,
	}

	q := lane.NewQueue()
	q.Enqueue(root)
	visited := make(map[string]bool)
	visited["/"] = true
	var crawled urlset

	for !q.Empty() {
		f := q.Dequeue().(node)
		l := f.Link.Href

		if f.Depth > depth {
			break
		}

		if visited[l] {
			continue
		}
		visited[l] = true
		fmt.Println("Depth:", f.Depth, "Link", l)
		crawled.URLs = append(crawled.URLs, url{
			Loc: l,
		})

		res, err := http.Get(l)
		if err != nil {
			continue
		}
		links, _ := link.Parse(res.Body)
		for _, c := range links {
			if !isSameDomain(c, domain) {
				continue
			}
			nl := normalizeLink(c, domain)
			if visited[nl] {
				continue
			}
			n := node{
				Link: link.Link{
					Href: nl,
					Text: c.Text,
				},
				Depth: f.Depth + 1,
			}
			q.Enqueue(n)
		}
	}

	return crawled
}

func isRelativeLink(link link.Link) bool {
	if strings.HasPrefix(link.Href, "/") && !strings.HasPrefix(link.Href, "//") {
		return true
	}
	return false
}

func isSameDomain(link link.Link, domain string) bool {
	if strings.HasPrefix(link.Href, domain) || isRelativeLink(link) {
		return true
	}
	return false
}

func normalizeLink(link link.Link, domain string) string {
	if isRelativeLink(link) {
		return domain + link.Href
	}
	return link.Href
}
