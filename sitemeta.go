package sitemeta

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

type Config struct {
	HTTPTimeout    time.Duration
	ChromeTimeout  time.Duration
	ChromeWaitTime time.Duration
	UserAgent      string
}

func DefaultConfig() *Config {
	return &Config{
		HTTPTimeout:    10 * time.Second,
		ChromeTimeout:  20 * time.Second,
		ChromeWaitTime: 1 * time.Second,
		UserAgent:      "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	}
}

type SiteMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
}

type Client struct {
	config *Config
	client *http.Client
}

func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		config: config,
		client: &http.Client{
			Timeout: config.HTTPTimeout,
		},
	}
}

func (c *Client) GetSiteMeta(websiteURL string) (*SiteMeta, error) {
	if websiteURL == "" {
		return nil, fmt.Errorf("website URL cannot be empty")
	}

	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	meta, err := c.extractMetaWithHTTP(parsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata with HTTP: %w", err)
	}

	if meta.Description == "" {
		chromedpMeta, err := c.extractMetaWithChromedp(parsedURL.String())
		if err != nil {
			fmt.Printf("ChromeDP extraction failed: %v, returning HTTP result\n", err)
			return meta, nil
		}
		return chromedpMeta, nil
	}

	return meta, nil
}

func GetSiteMeta(websiteURL string) (*SiteMeta, error) {
	client := NewClient(nil)
	return client.GetSiteMeta(websiteURL)
}

func (c *Client) extractMetaWithHTTP(websiteURL string) (*SiteMeta, error) {
	doc, err := c.fetchHTML(websiteURL)
	if err != nil {
		return nil, err
	}

	return c.parseSiteMeta(doc, websiteURL), nil
}

func (c *Client) extractMetaWithChromedp(websiteURL string) (*SiteMeta, error) {
	doc, err := c.renderDOMWithChrome(websiteURL)
	if err != nil {
		return nil, err
	}

	return c.parseSiteMeta(doc, websiteURL), nil
}

func (c *Client) fetchHTML(websiteURL string) (*html.Node, error) {
	req, err := http.NewRequest("GET", websiteURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.config.UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

func (c *Client) renderDOMWithChrome(target string) (*html.Node, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, c.config.ChromeTimeout)
	defer cancel()

	var htmlStr string
	err := chromedp.Run(ctx,
		chromedp.Navigate(target),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(c.config.ChromeWaitTime),
		chromedp.OuterHTML("html", &htmlStr, chromedp.ByQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("chrome rendering failed: %w", err)
	}

	doc, err := html.Parse(bytes.NewReader([]byte(htmlStr)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse rendered HTML: %w", err)
	}

	return doc, nil
}

func (c *Client) parseSiteMeta(doc *html.Node, siteURL string) *SiteMeta {
	meta := &SiteMeta{URL: siteURL}

	head := c.findHeadTag(doc)
	if head == nil {
		return meta
	}

	titleNode := c.findTitleTag(head)
	if titleNode != nil && titleNode.FirstChild != nil {
		meta.Title = titleNode.FirstChild.Data
	}

	metaTags := c.findMetaTags(head)
	meta.Description = c.extractDescription(metaTags)
	
	if img := c.extractImage(metaTags); img != "" {
		meta.Image = c.resolveImageURL(img, siteURL)
	}

	return meta
}

func (c *Client) findHeadTag(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "head" {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if head := c.findHeadTag(child); head != nil {
			return head
		}
	}

	return nil
}

func (c *Client) findTitleTag(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "title" {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if title := c.findTitleTag(child); title != nil {
			return title
		}
	}

	return nil
}

func (c *Client) findMetaTags(node *html.Node) []*html.Node {
	var metaTags []*html.Node

	if node.Type == html.ElementNode && node.Data == "meta" {
		metaTags = append(metaTags, node)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		metaTags = append(metaTags, c.findMetaTags(child)...)
	}

	return metaTags
}

func (c *Client) extractDescription(metaTags []*html.Node) string {
	for _, meta := range metaTags {
		if len(meta.Attr) < 2 {
			continue
		}

		if meta.Attr[0].Key == "name" && meta.Attr[0].Val == "description" {
			return meta.Attr[1].Val
		}

		if meta.Attr[0].Key == "property" && meta.Attr[0].Val == "og:description" {
			return meta.Attr[1].Val
		}

		if meta.Attr[0].Key == "name" && meta.Attr[0].Val == "twitter:description" {
			return meta.Attr[1].Val
		}
	}

	return ""
}

func (c *Client) extractImage(metaTags []*html.Node) string {
	for _, meta := range metaTags {
		if len(meta.Attr) < 2 {
			continue
		}

		if meta.Attr[0].Key == "property" && meta.Attr[0].Val == "og:image" {
			return meta.Attr[1].Val
		}

		if meta.Attr[0].Key == "name" && meta.Attr[0].Val == "twitter:image" {
			return meta.Attr[1].Val
		}
	}

	return ""
}

func (c *Client) resolveImageURL(imageURL, baseURL string) string {
	if imageURL == "" {
		return ""
	}

	if parsed, err := url.Parse(imageURL); err == nil && parsed.IsAbs() {
		return imageURL
	}

	if base, err := url.Parse(baseURL); err == nil {
		if resolved, err := base.Parse(imageURL); err == nil {
			return resolved.String()
		}
	}

	return imageURL
}