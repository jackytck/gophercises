package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jackytck/gophercises/ex4/link"
	lane "gopkg.in/oleiade/lane.v1"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func main() {
	domain := flag.String("domain", "https://www.golang.org", "Root domain")
	depth := flag.Int("depth", 1, "Depth of transversal")
	out := flag.String("out", "sitemap.xml", "Output filename")
	flag.Parse()

	base, err := getBaseURL(*domain)
	if err != nil {
		panic(err)
	}

	urls := crawl(base, *depth)
	if err := generateXML(urls, *out); err != nil {
		panic(err)
	}
}

func getBaseURL(reqURL string) (string, error) {
	r, err := http.Get(reqURL)
	if err != nil {
		return "", err
	}
	u := r.Request.URL
	b := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}
	return b.String(), nil
}

func generateXML(urls urlset, out string) error {
	f, err := os.Create(out)
	defer f.Close()
	if err != nil {
		return err
	}
	urls.XMLns = xmlns
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
	URLs    []loc
}

type loc struct {
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
		crawled.URLs = append(crawled.URLs, loc{
			Loc: l,
		})

		res, err := http.Get(l)
		if err != nil {
			res.Body.Close()
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

		res.Body.Close()
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
