package goplaces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSearchSuccess(t *testing.T) {
	var gotRequest map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/places:searchText" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Goog-Api-Key") != "test-key" {
			t.Fatalf("missing api key header")
		}
		if r.Header.Get("X-Goog-FieldMask") != searchFieldMask {
			t.Fatalf("unexpected field mask: %s", r.Header.Get("X-Goog-FieldMask"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &gotRequest); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "places": [
    {
      "id": "abc",
      "displayName": {"text": "Cafe"},
      "formattedAddress": "123 Street",
      "location": {"latitude": 1.23, "longitude": 4.56},
      "rating": 4.7,
      "priceLevel": "PRICE_LEVEL_MODERATE",
      "types": ["cafe"],
      "currentOpeningHours": {"openNow": true}
    }
  ],
  "nextPageToken": "next"
}`))
	}))
	defer server.Close()

	client := NewClient(Options{
		APIKey:  "test-key",
		BaseURL: server.URL + "/v1",
		Timeout: time.Second,
	})

	open := true
	minRating := 4.0
	request := SearchRequest{
		Query:     "coffee",
		Limit:     5,
		PageToken: "token",
		Language:  "en",
		Region:    "US",
		Filters: &Filters{
			Keyword:     "best",
			Types:       []string{"cafe"},
			OpenNow:     &open,
			MinRating:   &minRating,
			PriceLevels: []int{2},
		},
		LocationBias: &LocationBias{Lat: 40.0, Lng: -70.0, RadiusM: 500},
	}

	response, err := client.Search(context.Background(), request)
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if len(response.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(response.Results))
	}
	result := response.Results[0]
	if result.PlaceID != "abc" {
		t.Fatalf("unexpected place id: %s", result.PlaceID)
	}
	if result.Name != "Cafe" {
		t.Fatalf("unexpected name: %s", result.Name)
	}
	if result.PriceLevel == nil || *result.PriceLevel != 2 {
		t.Fatalf("unexpected price level: %#v", result.PriceLevel)
	}
	if result.OpenNow == nil || *result.OpenNow != true {
		t.Fatalf("unexpected openNow: %#v", result.OpenNow)
	}
	if response.NextPageToken != "next" {
		t.Fatalf("unexpected token: %s", response.NextPageToken)
	}

	if gotRequest["textQuery"] != "coffee best" {
		t.Fatalf("unexpected textQuery: %#v", gotRequest["textQuery"])
	}
	if gotRequest["pageSize"].(float64) != 5 {
		t.Fatalf("unexpected pageSize: %#v", gotRequest["pageSize"])
	}
	if gotRequest["pageToken"] != "token" {
		t.Fatalf("unexpected pageToken: %#v", gotRequest["pageToken"])
	}
	if gotRequest["languageCode"] != "en" {
		t.Fatalf("unexpected languageCode: %#v", gotRequest["languageCode"])
	}
	if gotRequest["regionCode"] != "US" {
		t.Fatalf("unexpected regionCode: %#v", gotRequest["regionCode"])
	}
}

func TestSearchHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad"))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL})
	_, err := client.Search(context.Background(), SearchRequest{Query: "coffee"})
	var apiErr *APIError
	if err == nil || !errors.As(err, &apiErr) {
		t.Fatalf("expected api error, got %v", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}

func TestSearchInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL})
	_, err := client.Search(context.Background(), SearchRequest{Query: "coffee"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAutocompleteSuccess(t *testing.T) {
	var gotRequest map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/places:autocomplete" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Goog-FieldMask") != autocompleteFieldMask {
			t.Fatalf("unexpected field mask: %s", r.Header.Get("X-Goog-FieldMask"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &gotRequest); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		_, _ = w.Write([]byte(`{
  "suggestions": [
    {
      "placePrediction": {
        "placeId": "place-1",
        "text": {"text": "Coffee Bar"},
        "structuredFormat": {
          "mainText": {"text": "Coffee"},
          "secondaryText": {"text": "Seattle"}
        },
        "types": ["cafe"]
      }
    },
    {
      "queryPrediction": {
        "text": {"text": "coffee beans"},
        "structuredFormat": {
          "mainText": {"text": "coffee beans"},
          "secondaryText": {"text": "query"}
        }
      }
    }
  ]
}`))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL + "/v1"})
	response, err := client.Autocomplete(context.Background(), AutocompleteRequest{
		Input:        "cof",
		Limit:        5,
		SessionToken: "session",
		Language:     "en",
		Region:       "US",
		LocationBias: &LocationBias{Lat: 1.1, Lng: 2.2, RadiusM: 100},
	})
	if err != nil {
		t.Fatalf("autocomplete error: %v", err)
	}
	if len(response.Suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(response.Suggestions))
	}
	if response.Suggestions[0].Kind != "place" || response.Suggestions[0].PlaceID != "place-1" {
		t.Fatalf("unexpected place suggestion: %#v", response.Suggestions[0])
	}
	if response.Suggestions[1].Kind != "query" || response.Suggestions[1].Text != "coffee beans" {
		t.Fatalf("unexpected query suggestion: %#v", response.Suggestions[1])
	}

	if gotRequest["input"] != "cof" {
		t.Fatalf("unexpected input: %#v", gotRequest["input"])
	}
	if gotRequest["sessionToken"] != "session" {
		t.Fatalf("unexpected session token: %#v", gotRequest["sessionToken"])
	}
	if gotRequest["languageCode"] != "en" {
		t.Fatalf("unexpected languageCode: %#v", gotRequest["languageCode"])
	}
	if gotRequest["regionCode"] != "US" {
		t.Fatalf("unexpected regionCode: %#v", gotRequest["regionCode"])
	}
	locationBias := gotRequest["locationBias"].(map[string]any)
	if locationBias["circle"] == nil {
		t.Fatalf("missing location bias circle")
	}
}

func TestAutocompleteLimitTrims(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
  "suggestions": [
    {"queryPrediction": {"text": {"text": "a"}}},
    {"queryPrediction": {"text": {"text": "b"}}}
  ]
}`))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL})
	response, err := client.Autocomplete(context.Background(), AutocompleteRequest{
		Input: "cof",
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("autocomplete error: %v", err)
	}
	if len(response.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(response.Suggestions))
	}
}

func TestDetailsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/places/place-123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("languageCode") != "en" {
			t.Fatalf("unexpected languageCode: %s", r.URL.Query().Get("languageCode"))
		}
		if r.URL.Query().Get("regionCode") != "US" {
			t.Fatalf("unexpected regionCode: %s", r.URL.Query().Get("regionCode"))
		}
		if r.Header.Get("X-Goog-FieldMask") != detailsFieldMaskBase {
			t.Fatalf("unexpected field mask: %s", r.Header.Get("X-Goog-FieldMask"))
		}
		_, _ = w.Write([]byte(`{
  "id": "place-123",
  "displayName": {"text": "Park"},
  "formattedAddress": "Central",
  "location": {"latitude": 10, "longitude": 20},
  "rating": 4.2,
  "priceLevel": "PRICE_LEVEL_FREE",
  "types": ["park"],
  "regularOpeningHours": {"weekdayDescriptions": ["Mon: 9-5"]},
  "currentOpeningHours": {"openNow": false},
  "nationalPhoneNumber": "+1 555",
  "websiteUri": "https://example.com"
}`))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL + "/v1"})
	place, err := client.DetailsWithOptions(context.Background(), DetailsRequest{
		PlaceID:  "place-123",
		Language: "en",
		Region:   "US",
	})
	if err != nil {
		t.Fatalf("details error: %v", err)
	}
	if place.PlaceID != "place-123" {
		t.Fatalf("unexpected id: %s", place.PlaceID)
	}
	if place.OpenNow == nil || *place.OpenNow != false {
		t.Fatalf("unexpected openNow")
	}
	if len(place.Hours) != 1 {
		t.Fatalf("unexpected hours")
	}
}

func TestDetailsWithReviews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("X-Goog-FieldMask"), "reviews") {
			t.Fatalf("expected reviews in field mask: %s", r.Header.Get("X-Goog-FieldMask"))
		}
		_, _ = w.Write([]byte(`{
  "id": "place-123",
  "reviews": [
    {
      "name": "places/place-123/reviews/1",
      "rating": 4.5,
      "text": {"text": "Great coffee", "languageCode": "en"},
      "authorAttribution": {"displayName": "Alice", "uri": "https://example.com"},
      "relativePublishTimeDescription": "2 weeks ago",
      "publishTime": "2024-01-01T00:00:00Z",
      "visitDate": {"year": 2024, "month": 1, "day": 2}
    }
  ]
}`))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL + "/v1"})
	details, err := client.DetailsWithOptions(context.Background(), DetailsRequest{
		PlaceID:        "place-123",
		IncludeReviews: true,
	})
	if err != nil {
		t.Fatalf("details error: %v", err)
	}
	if len(details.Reviews) != 1 {
		t.Fatalf("expected 1 review")
	}
	review := details.Reviews[0]
	if review.Author == nil || review.Author.DisplayName != "Alice" {
		t.Fatalf("unexpected author: %#v", review.Author)
	}
	if review.Text == nil || review.Text.Text != "Great coffee" {
		t.Fatalf("unexpected text: %#v", review.Text)
	}
	if review.VisitDate == nil || review.VisitDate.Year != 2024 {
		t.Fatalf("unexpected visit date: %#v", review.VisitDate)
	}
}

func TestDetailsFieldMaskForRequest(t *testing.T) {
	req := DetailsRequest{}
	if got := detailsFieldMaskForRequest(req); got != detailsFieldMaskBase {
		t.Fatalf("unexpected field mask: %s", got)
	}
	req.IncludeReviews = true
	got := detailsFieldMaskForRequest(req)
	if !strings.Contains(got, "reviews") {
		t.Fatalf("expected reviews in field mask: %s", got)
	}
}

func TestResolveSuccess(t *testing.T) {
	var gotRequest map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Goog-FieldMask") != resolveFieldMask {
			t.Fatalf("unexpected field mask: %s", r.Header.Get("X-Goog-FieldMask"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &gotRequest); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		_, _ = w.Write([]byte(`{
  "places": [
    {
      "id": "loc-1",
      "displayName": {"text": "Downtown"},
      "formattedAddress": "Main",
      "location": {"latitude": 1, "longitude": 2},
      "types": ["neighborhood"]
    }
  ]
}`))
	}))
	defer server.Close()

	client := NewClient(Options{APIKey: "test-key", BaseURL: server.URL})
	response, err := client.Resolve(context.Background(), LocationResolveRequest{
		LocationText: "Downtown",
		Language:     "en",
		Region:       "US",
	})
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if len(response.Results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if gotRequest["languageCode"] != "en" {
		t.Fatalf("unexpected languageCode: %#v", gotRequest["languageCode"])
	}
	if gotRequest["regionCode"] != "US" {
		t.Fatalf("unexpected regionCode: %#v", gotRequest["regionCode"])
	}
}

func TestMissingAPIKey(t *testing.T) {
	client := NewClient(Options{})
	_, err := client.Search(context.Background(), SearchRequest{Query: "coffee"})
	if !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("expected missing api key error")
	}
}

func TestValidationErrors(t *testing.T) {
	client := NewClient(Options{APIKey: "test-key", BaseURL: "http://example.com"})

	_, err := client.Search(context.Background(), SearchRequest{Query: ""})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	minRating := 9.0
	_, err = client.Search(context.Background(), SearchRequest{Query: "coffee", Filters: &Filters{MinRating: &minRating}})
	if err == nil {
		t.Fatalf("expected rating error")
	}

	_, err = client.Search(context.Background(), SearchRequest{Query: "coffee", Limit: 42})
	if err == nil {
		t.Fatalf("expected limit error")
	}

	_, err = client.Search(context.Background(), SearchRequest{Query: "coffee", Filters: &Filters{PriceLevels: []int{9}}})
	if err == nil {
		t.Fatalf("expected price level error")
	}

	_, err = client.Search(context.Background(), SearchRequest{Query: "coffee", LocationBias: &LocationBias{Lat: 200, Lng: 0, RadiusM: 1}})
	if err == nil {
		t.Fatalf("expected location error")
	}

	_, err = client.Resolve(context.Background(), LocationResolveRequest{LocationText: ""})
	if err == nil {
		t.Fatalf("expected resolve error")
	}

	_, err = client.Resolve(context.Background(), LocationResolveRequest{LocationText: "x", Limit: 99})
	if err == nil {
		t.Fatalf("expected resolve limit error")
	}

	_, err = client.Autocomplete(context.Background(), AutocompleteRequest{Input: ""})
	if err == nil {
		t.Fatalf("expected autocomplete input error")
	}

	_, err = client.Autocomplete(context.Background(), AutocompleteRequest{Input: "x", Limit: 99})
	if err == nil {
		t.Fatalf("expected autocomplete limit error")
	}

	_, err = client.Details(context.Background(), "")
	if err == nil {
		t.Fatalf("expected details error")
	}
}

func TestBuildSearchBodyOmitsEmptyPriceLevels(t *testing.T) {
	request := SearchRequest{Query: "coffee", Filters: &Filters{PriceLevels: []int{9}}}
	body := buildSearchBody(request)
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if bytes.Contains(payload, []byte("priceLevels")) {
		t.Fatalf("unexpected priceLevels in payload")
	}
}

func TestMappingHelpers(t *testing.T) {
	if mapLatLng(nil) != nil {
		t.Fatalf("expected nil location")
	}
	if displayName(nil) != "" {
		t.Fatalf("expected empty display name")
	}
	if openNow(nil) != nil {
		t.Fatalf("expected nil open now")
	}
	if weekdayDescriptions(nil) != nil {
		t.Fatalf("expected nil hours")
	}
	if mapPriceLevel("UNKNOWN") != nil {
		t.Fatalf("expected nil price level")
	}
}
