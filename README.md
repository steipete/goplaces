# 📍 goplaces — Modern Places for Go

Modern Go client + CLI for the Google Places API (New). Fast for humans, tidy for scripts.

## Highlights

- Text search with filters: keyword, type, open now, min rating, price levels.
- Autocomplete suggestions for places + queries (session tokens supported).
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

### Getting a Google Places API Key

1. **Create a Google Cloud Project**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Click "Select a project" → "New Project"
   - Name it (e.g., "goplaces") and click "Create"

2. **Enable the Places API (New)**
   - Go to [APIs & Services → Library](https://console.cloud.google.com/apis/library)
   - Search for "Places API (New)" — make sure it says **(New)**!
   - Click "Enable"

3. **Create an API Key**
   - Go to [APIs & Services → Credentials](https://console.cloud.google.com/apis/credentials)
   - Click "Create Credentials" → "API Key"
   - Copy the key

4. **Set the Environment Variable**
   ```bash
   export GOOGLE_PLACES_API_KEY="your-api-key-here"
   ```
   Add to your `~/.zshrc` or `~/.bashrc` to persist.

5. **(Recommended) Restrict the Key**
   - Click on the key in Credentials
   - Under "API restrictions", select "Restrict key" → "Places API (New)"
   - Set quota limits in [Quotas](https://console.cloud.google.com/apis/api/places.googleapis.com/quotas)

> **Note**: The Places API has usage costs. Check [pricing](https://developers.google.com/maps/documentation/places/web-service/usage-and-billing) and set budget alerts!

## CLI

```text
goplaces [--api-key=KEY] [--base-url=URL] [--timeout=10s] [--json] [--no-color] [--verbose]
         <command>

Commands:
  autocomplete  Autocomplete places and queries.
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

Autocomplete:

```bash
goplaces autocomplete "cof" --session-token "goplaces-demo" --limit 5 --language en --region US
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

autocomplete, err := client.Autocomplete(ctx, goplaces.AutocompleteRequest{
    Input:        "cof",
    SessionToken: "goplaces-demo",
    Limit:        5,
    Language:     "en",
    Region:       "US",
})
```

## Notes

- `Filters.Types` maps to `includedType` (Google accepts a single value). Only the first type is sent.
- Price levels map to Google enums: `0` (free) → `4` (very expensive).
- Reviews are returned only when `IncludeReviews`/`--reviews` is set.
- Field masks are defined alongside each request (e.g. `search.go`, `details.go`, `autocomplete.go`).
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

- Use a custom endpoint (proxy/mock): `GOOGLE_PLACES_E2E_BASE_URL`
- Override the search text used in E2E: `GOOGLE_PLACES_E2E_QUERY`
- Override language code for E2E: `GOOGLE_PLACES_E2E_LANGUAGE`
- Override region code for E2E: `GOOGLE_PLACES_E2E_REGION`
