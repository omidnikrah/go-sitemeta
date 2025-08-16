# go-sitemeta

A clean and efficient Go package for extracting metadata from websites, including support for JavaScript-rendered content.

## Features

- **Simple HTTP extraction**: Fast metadata extraction for static websites
- **ChromeDP integration**: JavaScript-rendered content support when needed
- **Configurable**: Customizable timeouts, user agents, and settings
- **Error handling**: Comprehensive error handling with meaningful messages
- **Clean API**: Well-structured, documented, and easy to use

## Installation

```bash
go get github.com/omidnikrah/go-sitemeta
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/omidnikrah/go-sitemeta"
)

func main() {
    // Simple usage with default configuration
    meta, err := sitemeta.GetSiteMeta("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Title: %s\n", meta.Title)
    fmt.Printf("Description: %s\n", meta.Description)
    fmt.Printf("Image: %s\n", meta.Image)
    fmt.Printf("URL: %s\n", meta.URL)
}
```

### Advanced Usage with Custom Configuration

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/omidnikrah/go-sitemeta"
)

func main() {
    // Create custom configuration
    config := &sitemeta.Config{
        HTTPTimeout:    5 * time.Second,
        ChromeTimeout:  15 * time.Second,
        ChromeWaitTime: 2 * time.Second,
        UserAgent:      "MyBot/1.0",
    }
    
    // Create client with custom configuration
    client := sitemeta.NewClient(config)
    
    // Extract metadata
    meta, err := client.GetSiteMeta("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Title: %s\n", meta.Title)
    fmt.Printf("Description: %s\n", meta.Description)
    fmt.Printf("Image: %s\n", meta.Image)
    fmt.Printf("URL: %s\n", meta.URL)
}
```

## API Reference

### Types

#### `SiteMeta`

Represents the metadata extracted from a website.

```go
type SiteMeta struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Image       string `json:"image"`
    URL         string `json:"url"`
}
```

#### `Config`

Configuration options for the sitemeta client.

```go
type Config struct {
    HTTPTimeout    time.Duration // HTTP request timeout
    ChromeTimeout  time.Duration // ChromeDP rendering timeout
    ChromeWaitTime time.Duration // Wait time after page load
    UserAgent      string        // User agent string
}
```

#### `Client`

A sitemeta client with configuration.

```go
type Client struct {
    config *Config
    client *http.Client
}
```

## Supported Meta Tags

The package extracts metadata from the following meta tags:

### Title
- `<title>` tag

### Description
- `name="description"`
- `property="og:description"`
- `name="twitter:description"`

### Image
- `property="og:image"`
- `name="twitter:image"`

## Testing

Run the tests:

```bash
go test
```

## License

This project is licensed under the MIT License.
