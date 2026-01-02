package cli

import (
	"time"
)

// Root defines the CLI command tree.
type Root struct {
	Global       GlobalOptions   `embed:""`
	Autocomplete AutocompleteCmd `cmd:"" help:"Autocomplete places and queries."`
	Search       SearchCmd       `cmd:"" help:"Search places by text query."`
	Details      DetailsCmd      `cmd:"" help:"Fetch place details by place ID."`
	Resolve      ResolveCmd      `cmd:"" help:"Resolve a location string to candidate places."`
}

// GlobalOptions are flags shared by all commands.
type GlobalOptions struct {
	APIKey  string        `help:"Google Places API key." env:"GOOGLE_PLACES_API_KEY"`
	BaseURL string        `help:"Places API base URL." env:"GOOGLE_PLACES_BASE_URL" default:"https://places.googleapis.com/v1"`
	Timeout time.Duration `help:"HTTP timeout." default:"10s"`
	JSON    bool          `help:"Output JSON."`
	NoColor bool          `help:"Disable color output."`
	Verbose bool          `help:"Verbose logging."`
	Version VersionFlag   `name:"version" help:"Print version and exit."`
}

// SearchCmd runs text search queries.
type SearchCmd struct {
	Query      string   `arg:"" name:"query" help:"Search text."`
	Limit      int      `help:"Max results (1-20)." default:"10"`
	PageToken  string   `help:"Page token for pagination."`
	Language   string   `help:"BCP-47 language code (e.g. en, en-US)."`
	Region     string   `help:"CLDR region code (e.g. US, DE)."`
	Keyword    string   `help:"Keyword to append to the query."`
	Type       []string `help:"Place type filter (includedType). Repeatable."`
	OpenNow    *bool    `help:"Return only currently open places."`
	MinRating  *float64 `help:"Minimum rating (0-5)."`
	PriceLevel []int    `help:"Price levels 0-4. Repeatable."`
	Lat        *float64 `help:"Latitude for location bias."`
	Lng        *float64 `help:"Longitude for location bias."`
	RadiusM    *float64 `help:"Radius in meters for location bias."`
}

// AutocompleteCmd runs autocomplete queries.
type AutocompleteCmd struct {
	Input        string   `arg:"" name:"input" help:"Autocomplete input text."`
	Limit        int      `help:"Max suggestions (1-20)." default:"5"`
	SessionToken string   `help:"Session token for billing consistency."`
	Language     string   `help:"BCP-47 language code (e.g. en, en-US)."`
	Region       string   `help:"CLDR region code (e.g. US, DE)."`
	Lat          *float64 `help:"Latitude for location bias."`
	Lng          *float64 `help:"Longitude for location bias."`
	RadiusM      *float64 `help:"Radius in meters for location bias."`
}

// DetailsCmd fetches place details.
type DetailsCmd struct {
	PlaceID  string `arg:"" name:"place_id" help:"Place ID."`
	Language string `help:"BCP-47 language code (e.g. en, en-US)."`
	Region   string `help:"CLDR region code (e.g. US, DE)."`
	Reviews  bool   `help:"Include reviews in the response."`
}

// ResolveCmd resolves a location string into candidates.
type ResolveCmd struct {
	LocationText string `arg:"" name:"location" help:"Location text to resolve."`
	Limit        int    `help:"Max results (1-10)." default:"5"`
	Language     string `help:"BCP-47 language code (e.g. en, en-US)."`
	Region       string `help:"CLDR region code (e.g. US, DE)."`
}
