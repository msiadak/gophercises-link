package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	HREF string
	Text string
}

func ExtractLinks(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	linkCh := make(chan Link)
	done := make(chan bool)

	go nodeDFS(node, linkCh, 0)

	links := make([]Link, 0)
	for link := range linkCh {
		links = append(links, link)
	}
	return links, nil
}

func nodeDFS(n *html.Node, links chan Link, depth int) {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				fragments := make([]string, 0)
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseInnerHTML(c, &fragments)
				}
				link := Link{attr.Val, strings.Join(fragments, " ")}
				*links = append(*links, link)
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodeDFS(c, links, depth+1)
	}
	done <- true
}

func extractInnerText(n *html.Node, fragments *[]string) {
	if n.Type == html.TextNode {
		*fragments = append(*fragments, strings.TrimSpace(n.Data))
	}
	if n.Type == html.ElementNode {
		extractInnerText(n, fragments)
	}
}
