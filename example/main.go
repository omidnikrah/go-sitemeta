package main

import (
	"fmt"
	"log"
	"time"

	"github.com/omidnikrah/go-sitemeta"
)

func main() {
	// Example 1: Simple usage with default configuration
	fmt.Println("=== Example 1: Simple Usage ===")
	meta, err := sitemeta.GetSiteMeta("https://omid.toys")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printMeta(meta)
	}

	fmt.Println("\n=== Example 2: Custom Configuration ===")
	// Example 2: Custom configuration
	config := &sitemeta.Config{
		HTTPTimeout:    5 * time.Second,
		ChromeTimeout:  15 * time.Second,
		ChromeWaitTime: 2 * time.Second,
		UserAgent:      "MyCustomBot/1.0",
	}

	client := sitemeta.NewClient(config)
	meta, err = client.GetSiteMeta("https://omid.toys")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printMeta(meta)
	}

}

func printMeta(meta *sitemeta.SiteMeta) {
	fmt.Printf("Title: %s\n", meta.Title)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Printf("Image: %s\n", meta.Image)
	fmt.Printf("URL: %s\n", meta.URL)
}
