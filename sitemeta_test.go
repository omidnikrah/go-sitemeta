package sitemeta

import (
	"testing"
	"time"
)

func TestGetSiteMeta(t *testing.T) {
	url := "https://omid.toys"
	meta, err := GetSiteMeta(url)
	
	if err != nil {
		t.Errorf("GetSiteMeta failed: %v", err)
	}
	
	if meta.URL != url {
		t.Errorf("Expected URL %s, got %s", url, meta.URL)
	}
	
	// Test with invalid URL
	_, err = GetSiteMeta("http//omid.toys")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestNewClient(t *testing.T) {
	// Test with nil config
	client := NewClient(nil)
	if client == nil {
		t.Error("NewClient with nil config should return a client")
	}
	
	// Test with custom config
	config := &Config{
		HTTPTimeout:    5 * time.Second,
		ChromeTimeout:  10 * time.Second,
		ChromeWaitTime: 500 * time.Millisecond,
		UserAgent:      "Custom User Agent",
	}
	
	client = NewClient(config)
	if client == nil {
		t.Error("NewClient with custom config should return a client")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.HTTPTimeout == 0 {
		t.Error("DefaultConfig should set HTTPTimeout")
	}
	
	if config.ChromeTimeout == 0 {
		t.Error("DefaultConfig should set ChromeTimeout")
	}
	
	if config.ChromeWaitTime == 0 {
		t.Error("DefaultConfig should set ChromeWaitTime")
	}
	
	if config.UserAgent == "" {
		t.Error("DefaultConfig should set UserAgent")
	}
}

func TestClientGetSiteMeta(t *testing.T) {
	client := NewClient(nil)
	
	// Test with empty URL
	_, err := client.GetSiteMeta("")
	if err == nil {
		t.Error("Expected error for empty URL")
	}
	
	// Test with invalid URL
	_, err = client.GetSiteMeta("http//omid.toys")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestSiteMetaStruct(t *testing.T) {
	meta := &SiteMeta{
		Title:       "Test Title",
		Description: "Test Description",
		Image:       "https://omid.toys/og.png",
		URL:         "https://omid.toys",
	}
	
	if meta.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", meta.Title)
	}
	
	if meta.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got '%s'", meta.Description)
	}
	
	if meta.Image != "https://omid.toys/og.png" {
		t.Errorf("Expected image 'https://omid.toys/og.png', got '%s'", meta.Image)
	}
	
	if meta.URL != "https://omid.toys" {
		t.Errorf("Expected URL 'https://omid.toys', got '%s'", meta.URL)
	}
}