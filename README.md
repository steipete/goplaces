# 📍 goplaces

*Find places, Go fast*


Modern Go client + CLI for the Google Places API (New).

## Features

- Text search with filters, location bias, and pagination.
- Place details by place ID.
- Resolve a free-form location string into candidates.
- Human-readable output with color, plus `--json` for scripts.
- Strict validation + typed responses.

## Install

```bash
go get github.com/steipete/goplaces
```

## CLI

Set your API key via env or flag:

```bash
export GOOGLE_PLACES_API_KEY="..."
```

Search:

```bash
goplaces search "coffee" --min-rating 4 --open-now --limit 5 \
  --lat 40.8065 --lng -73.9719 --radius-m 3000
```

Details:

```bash
goplaces details ChIJN1t_tDeuEmsRUsoyG83frY4
```

Resolve:

```bash
goplaces resolve "Riverside Park, New York" --limit 5
```

JSON output:

```bash
goplaces search "pizza" --json
```

### CLI reference

```text
goplaces [--api-key=KEY] [--base-url=URL] [--timeout=10s] [--json] [--no-color] [--verbose]
         <command>

Commands:
  search   Search places by text query.
  details  Fetch place details by place ID.
  resolve  Resolve a location string to candidate places.
```

## Library

```go
boolPtr := func(v bool) *bool { return &v }
floatPtr := func(v float64) *float64 { return &v }

client := goplaces.NewClient(goplaces.Options{
    APIKey: os.Getenv("GOOGLE_PLACES_API_KEY"),
})

resp, err := client.Search(ctx, goplaces.SearchRequest{
    Query: "italian restaurant",
    Filters: &goplaces.Filters{
        OpenNow:   boolPtr(true),
        MinRating: floatPtr(4.0),
        Types:     []string{"restaurant"},
    },
    LocationBias: &goplaces.LocationBias{Lat: 40.8065, Lng: -73.9719, RadiusM: 3000},
    Limit: 10,
})
```

## Notes

- `Filters.Types` maps to `includedType` (Google supports a single value). Only the first type is sent.
- Price levels map to Google enums: `0` (free) → `4` (very expensive).
- Use `GOOGLE_PLACES_BASE_URL` to override the endpoint (useful for tests).

## Testing

```bash
go test ./... -coverprofile=coverage.out
```

Enforce coverage:

```bash
./scripts/check-coverage.sh
```
