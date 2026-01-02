package goplaces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	detailsFieldMaskBase   = "id,displayName,formattedAddress,location,rating,priceLevel,types,regularOpeningHours,currentOpeningHours,nationalPhoneNumber,websiteUri"
	detailsFieldMaskReview = "reviews"
)

// Details fetches details for a specific place ID.
func (c *Client) Details(ctx context.Context, placeID string) (PlaceDetails, error) {
	return c.DetailsWithOptions(ctx, DetailsRequest{PlaceID: placeID})
}

// DetailsWithOptions fetches place details with locale hints.
func (c *Client) DetailsWithOptions(ctx context.Context, req DetailsRequest) (PlaceDetails, error) {
	placeID := strings.TrimSpace(req.PlaceID)
	if placeID == "" {
		return PlaceDetails{}, ValidationError{Field: "place_id", Message: "required"}
	}

	endpoint, err := c.buildURL("/places/"+placeID, map[string]string{
		"languageCode": strings.TrimSpace(req.Language),
		"regionCode":   strings.TrimSpace(req.Region),
	})
	if err != nil {
		return PlaceDetails{}, err
	}

	payload, err := c.doRequest(ctx, http.MethodGet, endpoint, nil, detailsFieldMaskForRequest(req))
	if err != nil {
		return PlaceDetails{}, err
	}

	var place placeItem
	if err := json.Unmarshal(payload, &place); err != nil {
		return PlaceDetails{}, fmt.Errorf("goplaces: decode place details: %w", err)
	}

	return mapPlaceDetails(place), nil
}

func detailsFieldMaskForRequest(req DetailsRequest) string {
	if req.IncludeReviews {
		// Reviews are heavy; opt-in to include them.
		return detailsFieldMaskBase + "," + detailsFieldMaskReview
	}
	return detailsFieldMaskBase
}

func mapPlaceDetails(place placeItem) PlaceDetails {
	return PlaceDetails{
		PlaceID:    place.ID,
		Name:       displayName(place.DisplayName),
		Address:    place.FormattedAddress,
		Location:   mapLatLng(place.Location),
		Rating:     place.Rating,
		PriceLevel: mapPriceLevel(place.PriceLevel),
		Types:      place.Types,
		Phone:      place.NationalPhoneNumber,
		Website:    place.WebsiteURI,
		Hours:      weekdayDescriptions(place.RegularOpeningHours),
		OpenNow:    openNow(place.CurrentOpeningHours),
		Reviews:    mapReviews(place.Reviews),
	}
}
