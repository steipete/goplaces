# Changelog

## 0.1.0 - 2026-01-02

- Go client for Google Places API (New).
- Text search with filters: keyword, type, open now, min rating, price levels.
- Location bias (lat/lng/radius) and pagination tokens.
- Place details with hours, phone, website, rating, price, types.
- Optional reviews in details (`IncludeReviews` / `--reviews`).
- Resolve free-form location strings to candidate places.
- Locale hints (language + region) for search/resolve/details.
- Typed models, validation errors, and API error surfacing.
- CLI commands: `search`, `details`, `resolve` with color output + `--json`.
- Env/flag config: API key, base URL, timeouts, verbose logging.
- Lint/format guardrails + >= 90% coverage gate.
- GitHub Actions CI for tests, coverage, and linting.
