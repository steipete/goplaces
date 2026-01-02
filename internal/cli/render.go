package cli

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/steipete/goplaces"
)

func renderSearch(color Color, response goplaces.SearchResponse) string {
	var out bytes.Buffer
	count := len(response.Results)
	if count == 0 {
		return "No results."
	}
	out.WriteString(color.Bold(fmt.Sprintf("Results (%d)", count)))
	out.WriteString("\n")

	for i, place := range response.Results {
		out.WriteString(fmt.Sprintf("%d. %s\n", i+1, formatTitle(color, place.Name, place.Address)))
		writePlaceSummary(&out, color, place)
		if i < count-1 {
			out.WriteString("\n")
		}
	}

	if strings.TrimSpace(response.NextPageToken) != "" {
		out.WriteString("\n")
		out.WriteString(color.Dim("Next page token:"))
		out.WriteString(" ")
		out.WriteString(response.NextPageToken)
	}

	return out.String()
}

func renderDetails(color Color, place goplaces.PlaceDetails) string {
	var out bytes.Buffer
	out.WriteString(color.Bold(formatTitle(color, place.Name, place.Address)))
	out.WriteString("\n")
	writePlaceDetails(&out, color, place)
	return out.String()
}

func renderResolve(color Color, response goplaces.LocationResolveResponse) string {
	var out bytes.Buffer
	count := len(response.Results)
	if count == 0 {
		return "No results."
	}
	out.WriteString(color.Bold(fmt.Sprintf("Resolved (%d)", count)))
	out.WriteString("\n")

	for i, place := range response.Results {
		out.WriteString(fmt.Sprintf("%d. %s\n", i+1, formatTitle(color, place.Name, place.Address)))
		writeResolvedLocation(&out, color, place)
		if i < count-1 {
			out.WriteString("\n")
		}
	}
	return out.String()
}

func formatTitle(color Color, name string, address string) string {
	display := strings.TrimSpace(name)
	if display == "" {
		display = "(no name)"
	}
	if address == "" {
		return color.Cyan(display)
	}
	return color.Cyan(display) + " — " + address
}

func writePlaceSummary(out *bytes.Buffer, color Color, place goplaces.PlaceSummary) {
	writeLine(out, color, "ID", place.PlaceID)
	writeLocation(out, color, place.Location)
	writeRating(out, color, place.Rating, place.PriceLevel)
	writeTypes(out, color, place.Types)
	writeOpenNow(out, color, place.OpenNow)
}

func writePlaceDetails(out *bytes.Buffer, color Color, place goplaces.PlaceDetails) {
	writeLine(out, color, "ID", place.PlaceID)
	writeLocation(out, color, place.Location)
	writeRating(out, color, place.Rating, place.PriceLevel)
	writeTypes(out, color, place.Types)
	writeOpenNow(out, color, place.OpenNow)
	writeLine(out, color, "Phone", place.Phone)
	writeLine(out, color, "Website", place.Website)
	writeReviews(out, color, place.Reviews)
	if len(place.Hours) > 0 {
		out.WriteString(color.Dim("Hours:"))
		out.WriteString("\n")
		for _, entry := range place.Hours {
			out.WriteString("  - ")
			out.WriteString(entry)
			out.WriteString("\n")
		}
	}
}

func writeResolvedLocation(out *bytes.Buffer, color Color, place goplaces.ResolvedLocation) {
	writeLine(out, color, "ID", place.PlaceID)
	writeLocation(out, color, place.Location)
	writeTypes(out, color, place.Types)
}

func writeReviews(out *bytes.Buffer, color Color, reviews []goplaces.Review) {
	if len(reviews) == 0 {
		return
	}
	out.WriteString(color.Dim("Reviews:"))
	out.WriteString("\n")

	// Keep CLI output compact by default.
	const maxReviews = 3
	count := len(reviews)
	limit := count
	if count > maxReviews {
		limit = maxReviews
	}

	for i := 0; i < limit; i++ {
		review := reviews[i]
		line := reviewLine(review)
		if line == "" {
			continue
		}
		out.WriteString("  - ")
		out.WriteString(line)
		out.WriteString("\n")
	}

	if count > maxReviews {
		out.WriteString(color.Dim(fmt.Sprintf("  ... %d more", count-maxReviews)))
		out.WriteString("\n")
	}
}

func writeLocation(out *bytes.Buffer, color Color, loc *goplaces.LatLng) {
	if loc == nil {
		return
	}
	writeLine(out, color, "Location", fmt.Sprintf("%.6f, %.6f", loc.Lat, loc.Lng))
}

func writeRating(out *bytes.Buffer, color Color, rating *float64, priceLevel *int) {
	if rating == nil && priceLevel == nil {
		return
	}
	parts := make([]string, 0, 2)
	if rating != nil {
		parts = append(parts, fmt.Sprintf("%.1f", *rating))
	}
	if priceLevel != nil {
		parts = append(parts, fmt.Sprintf("$%d", *priceLevel))
	}
	writeLine(out, color, "Rating", strings.Join(parts, " · "))
}

func writeTypes(out *bytes.Buffer, color Color, types []string) {
	if len(types) == 0 {
		return
	}
	unique := uniqueStrings(types)
	writeLine(out, color, "Types", strings.Join(unique, ", "))
}

func writeOpenNow(out *bytes.Buffer, color Color, openNow *bool) {
	if openNow == nil {
		return
	}
	value := "no"
	if *openNow {
		value = "yes"
	}
	writeLine(out, color, "Open now", value)
}

func writeLine(out *bytes.Buffer, color Color, label string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	out.WriteString(color.Dim(label + ":"))
	out.WriteString(" ")
	out.WriteString(value)
	out.WriteString("\n")
}

func reviewLine(review goplaces.Review) string {
	parts := make([]string, 0, 3)
	if review.Rating != nil {
		parts = append(parts, fmt.Sprintf("%.1f stars", *review.Rating))
	}
	if review.Author != nil && strings.TrimSpace(review.Author.DisplayName) != "" {
		parts = append(parts, "by "+review.Author.DisplayName)
	}
	if strings.TrimSpace(review.RelativePublishTimeDescription) != "" {
		parts = append(parts, "("+review.RelativePublishTimeDescription+")")
	}
	text := reviewText(review)
	if text != "" {
		parts = append(parts, text)
	}
	return strings.Join(parts, " ")
}

func reviewText(review goplaces.Review) string {
	text := ""
	if review.Text != nil {
		text = review.Text.Text
	}
	// Fall back to original text when translation is empty.
	if strings.TrimSpace(text) == "" && review.OriginalText != nil {
		text = review.OriginalText.Text
	}
	return truncateText(strings.TrimSpace(text), 200)
}

func truncateText(value string, maxLen int) string {
	if maxLen <= 0 || value == "" {
		return value
	}
	if len(value) <= maxLen {
		return value
	}
	// Byte-based truncation is OK here because we only display previews.
	return strings.TrimSpace(value[:maxLen]) + "..."
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
