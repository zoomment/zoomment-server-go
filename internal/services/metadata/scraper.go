package metadata

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// FetchSiteToken fetches the zoomment meta tag content from a URL
// This is equivalent to your fetchSiteToken function in Node.js
// It looks for: <meta name="zoomment" content="USER_ID">
func FetchSiteToken(siteURL string) (string, error) {
	// Make HTTP request
	resp, err := http.Get(siteURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse HTML and find meta tag
	return findZoommentMetaTag(string(body))
}

// findZoommentMetaTag parses HTML and finds the zoomment meta tag
func findZoommentMetaTag(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var token string
	var findMeta func(*html.Node)

	findMeta = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string

			for _, attr := range n.Attr {
				if attr.Key == "name" && strings.HasPrefix(attr.Val, "zoomment") {
					name = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}

			if name != "" && content != "" {
				token = content
				return
			}
		}

		// Recursively search child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findMeta(c)
			if token != "" {
				return
			}
		}
	}

	findMeta(doc)
	return token, nil
}

