package rss

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// httpClient is a reusable HTTP client
var httpClient = &http.Client{}

// RSS defines the RSS root structure
type RSS struct {
	Channel Channel `xml:"channel"`
}

// Channel contains metadata and a list of items
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents a single RSS feed item
type Item struct {
	Title       string     `xml:"title"`
	Links       []string   `xml:"link"` // Handles multiple <link> tags
	Description string     `xml:"description"`
	PubDate     string     `xml:"pubDate"`
	Guid        string     `xml:"guid"`
	Categories  []Category `xml:"category"`
}

// GenId generates a unique ID for the FeedItem based on its link and publication date
func (i *Item) GenId() string {

	input := ""

	if i.Guid != "" {
		input = i.Guid + i.Links[0] + i.PubDate
	} else {
		input = i.Links[0] + i.PubDate
	}

	// Generate a SHA-1 hash of the input string
	hash := sha1.Sum([]byte(input)) // returns [20]byte
	return hex.EncodeToString(hash[:])
}

// Category represents a category in an RSS item
type Category struct {
	Domain string `xml:"domain,attr"`
	Value  string `xml:",chardata"`
}

// GetLink returns the primary link from the item
func (i Item) GetLink() string {
	if len(i.Links) > 0 {
		return i.Links[0]
	}
	return ""
}

// HasCategory returns true if the item contains the given category value
func (i Item) HasCategory(keyword string) bool {
	for _, c := range i.Categories {
		if c.Value == keyword {
			return true
		}
	}
	return false
}

// AllCategoryDomains returns a list of unique category domains
func (i Item) AllCategoryDomains() []string {
	domains := make(map[string]struct{})
	for _, c := range i.Categories {
		domains[c.Domain] = struct{}{}
	}
	var out []string
	for d := range domains {
		out = append(out, d)
	}
	return out
}

// AllLinks returns all primary links from items
func (c Channel) AllLinks() []string {
	var out []string
	for _, item := range c.Items {
		if link := item.GetLink(); link != "" {
			out = append(out, link)
		}
	}
	return out
}

// FetchRSS downloads and parses RSS XML from the given URL
func FetchRSS(url string, ctx context.Context) (*RSS, error) {

	// Create a new HTTP request with a context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to create request: %w", err)
	}

	// Make the HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	// Unmarshal the XML data into the RSS struct
	rss := RSS{}
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, fmt.Errorf("unmarshal RSS xml failed: %w", err)
	}

	return &rss, nil
}
