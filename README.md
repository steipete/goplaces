# 📍 goplaces — Modern Places for Go

Modern Go client + CLI for the Google Places API (New). Fast for humans, tidy for scripts.

## Highlights

- Text search with filters: keyword, type, open now, min rating, price levels.
- Location bias (lat/lng/radius) and pagination tokens.
- Place details: hours, phone, website, rating, price, types.
- Optional reviews in details (`--reviews` / `IncludeReviews`).
- Resolve free-form location strings to candidate places.
- Locale hints (language + region) across search/resolve/details.
- Typed models, validation errors, and API error surfacing.
- CLI with color human output + `--json` (respects `NO_COLOR`).

## Install / Run

- Homebrew: `brew install steipete/tap/goplaces`
- Go: `go install github.com/steipete/goplaces/cmd/goplaces@latest`
- Source: `make goplaces`

## Config

```bash
export GOOGLE_PLACES_API_KEY="..."
```

Optional overrides:

- `GOOGLE_PLACES_BASE_URL` (testing, proxying, or mock servers)

### Getting an API key

1) Create a Google Cloud project.
2) Enable **Places API (New)** in the API Library.
3) Create an API key in **APIs & Services → Credentials**.
4) Restrict the key (HTTP referrers or IPs) and set quota/billing limits.

## CLI

```text
goplaces [--api-key=KEY] [--base-url=URL] [--timeout=10s] [--json] [--no-color] [--verbose]
         <command>

Commands:
  search   Search places by text query.
  details  Fetch place details by place ID.
  resolve  Resolve a location string to candidate places.
```

Search with filters + location bias:

```bash
goplaces search "coffee" --min-rating 4 --open-now --limit 5 \
  --lat 40.8065 --lng -73.9719 --radius-m 3000 --language en --region US
```

Pagination:

```bash
goplaces search "pizza" --page-token "NEXT_PAGE_TOKEN"
```

Details (with reviews):

```bash
goplaces details ChIJN1t_tDeuEmsRUsoyG83frY4 --reviews
```

Resolve:

```bash
goplaces resolve "Riverside Park, New York" --limit 5
```

JSON output:

```bash
goplaces search "sushi" --json
```

## Library

```go
boolPtr := func(v bool) *bool { return &v }
floatPtr := func(v float64) *float64 { return &v }

client := goplaces.NewClient(goplaces.Options{
    APIKey:  os.Getenv("GOOGLE_PLACES_API_KEY"),
    Timeout: 8 * time.Second,
})

search, err := client.Search(ctx, goplaces.SearchRequest{
    Query: "italian restaurant",
    Filters: &goplaces.Filters{
        OpenNow:   boolPtr(true),
        MinRating: floatPtr(4.0),
        Types:     []string{"restaurant"},
    },
    LocationBias: &goplaces.LocationBias{Lat: 40.8065, Lng: -73.9719, RadiusM: 3000},
    Language:     "en",
    Region:       "US",
    Limit:        10,
})

details, err := client.DetailsWithOptions(ctx, goplaces.DetailsRequest{
    PlaceID:        "ChIJN1t_tDeuEmsRUsoyG83frY4",
    Language:       "en",
    Region:         "US",
    IncludeReviews: true,
})
```

## Notes

- `Filters.Types` maps to `includedType` (Google accepts a single value). Only the first type is sent.
- Price levels map to Google enums: `0` (free) → `4` (very expensive).
- Reviews are returned only when `IncludeReviews`/`--reviews` is set.
- Field masks are defined in `client.go`; extend them if you need more fields.
- The Places API is billed and quota-limited; keep an eye on your Cloud Console quotas.

## Testing

```bash
make lint test coverage
```

### E2E tests (optional)

```bash
export GOOGLE_PLACES_API_KEY="..."
make e2e
```

Optional env overrides:

- `GOOGLE_PLACES_E2E_BASE_URL`
- `GOOGLE_PLACES_E2E_QUERY`
- `GOOGLE_PLACES_E2E_LANGUAGE`
- `GOOGLE_PLACES_E2E_REGION`
