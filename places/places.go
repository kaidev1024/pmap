package places

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kaidev1024/pugo/puslice"
	"googlemaps.github.io/maps"
)

var googlePlacesClient *maps.Client
var err error
var pointTypes map[string]struct{}

func Init(token string) error {
	googlePlacesClient, err = maps.NewClient(maps.WithAPIKey(token))
	if err != nil {
		return err
	}
	pointTypes = make(map[string]struct{})
	pointTypes["street_address"] = struct{}{}
	pointTypes["point_of_interest"] = struct{}{}
	pointTypes["intersection"] = struct{}{}
	pointTypes["premise"] = struct{}{}
	pointTypes["subpremise"] = struct{}{}
	pointTypes["establishment"] = struct{}{}
	return nil
}

func GetAutocompletePredictions(searchText, sessionToken string, isCity bool) ([]maps.AutocompletePrediction, error) {
	if isCity {
		return getAutocompletePredictionsForCity(searchText, sessionToken)
	}
	return getAutocompletePredictionsForPoint(searchText, sessionToken)
}

func GetDetailByPlaceID(placeID, sessionToken string) (maps.PlaceDetailsResult, error) {
	req := &maps.PlaceDetailsRequest{
		PlaceID:      placeID,
		SessionToken: parseSessionToken(sessionToken),
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskPlaceID,
			maps.PlaceDetailsFieldMaskName,
			maps.PlaceDetailsFieldMaskFormattedAddress,
			maps.PlaceDetailsFieldMaskGeometry,
			maps.PlaceDetailsFieldMaskAddressComponent,
		},
	}
	return googlePlacesClient.PlaceDetails(context.Background(), req)
}

func SearchByText(searchText string) (maps.PlacesSearchResult, error) {
	result, err := googlePlacesClient.TextSearch(context.Background(), &maps.TextSearchRequest{
		Query: searchText,
	})
	if err != nil {
		return maps.PlacesSearchResult{}, err
	}
	if len(result.Results) > 0 {
		return result.Results[0], nil
	}
	return maps.PlacesSearchResult{}, err
}

func parseSessionToken(sessionToken string) maps.PlaceAutocompleteSessionToken {
	sessionTokenUUID, err := uuid.Parse(sessionToken)
	if err != nil {
		log.Fatalf("invalid UUID: %v", err)
	}
	return maps.PlaceAutocompleteSessionToken(sessionTokenUUID)
}

func getAutocompletePredictionsForPoint(searchText, sessionToken string) ([]maps.AutocompletePrediction, error) {
	placesResp, err := googlePlacesClient.PlaceAutocomplete(context.Background(), &maps.PlaceAutocompleteRequest{
		Input:        searchText,
		SessionToken: parseSessionToken(sessionToken),
	})
	if err != nil {
		return nil, err
	}
	return puslice.Filter(placesResp.Predictions, func(prediction maps.AutocompletePrediction) bool {
		for _, t := range prediction.Types {
			if _, okay := pointTypes[t]; okay {
				return true
			}
		}
		return false
	}), nil
}

func getAutocompletePredictionsForCity(searchText, sessionToken string) ([]maps.AutocompletePrediction, error) {
	placesResp, err := googlePlacesClient.PlaceAutocomplete(context.Background(), &maps.PlaceAutocompleteRequest{
		Input:        searchText,
		SessionToken: parseSessionToken(sessionToken),
		Types:        maps.AutocompletePlaceTypeCities,
	})
	if err != nil {
		return nil, err
	}
	return placesResp.Predictions, nil
}
