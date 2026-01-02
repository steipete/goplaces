//go:build e2e
// +build e2e

package goplaces

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestE2ESearchAndDetails(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		t.Skip("GOOGLE_PLACES_API_KEY not set")
	}

	query := os.Getenv("GOOGLE_PLACES_E2E_QUERY")
	if query == "" {
		query = "coffee in Seattle"
	}
	language := os.Getenv("GOOGLE_PLACES_E2E_LANGUAGE")
	if language == "" {
		language = "en"
	}
	region := os.Getenv("GOOGLE_PLACES_E2E_REGION")
	if region == "" {
		region = "US"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := NewClient(Options{
		APIKey:  apiKey,
		BaseURL: os.Getenv("GOOGLE_PLACES_E2E_BASE_URL"),
		Timeout: 10 * time.Second,
	})

	search, err := client.Search(ctx, SearchRequest{
		Query:    query,
		Limit:    1,
		Language: language,
		Region:   region,
	})
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if len(search.Results) == 0 {
		t.Fatalf("expected search results")
	}

	placeID := search.Results[0].PlaceID
	if placeID == "" {
		t.Fatalf("expected place id")
	}

	details, err := client.DetailsWithOptions(ctx, DetailsRequest{
		PlaceID:  placeID,
		Language: language,
		Region:   region,
	})
	if err != nil {
		t.Fatalf("details error: %v", err)
	}
	if details.PlaceID == "" {
		t.Fatalf("expected details place id")
	}
}
