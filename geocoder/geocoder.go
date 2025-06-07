package geocoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aidenappl/RootedGeocoder/env"
	"github.com/aidenappl/RootedGeocoder/structs"
)

func GetGeocodedAddress(organisation structs.Organisation) (*structs.MapboxResponse, error) {

	baseURL := "https://api.mapbox.com/search/geocode/v6/forward"
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}

	params := url.Values{}
	params.Add("q", fmt.Sprintf("%s, %s, %s, %s", organisation.AddressLine1, organisation.City, organisation.State, organisation.ZipCode))
	params.Add("access_token", env.MapboxAccessToken) // Assuming you have a Mapbox access token in your env package

	parsedURL.RawQuery = params.Encode()

	// 2. Create HTTP Client and Request
	client := &http.Client{}
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// 3. Make the Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	var mapboxResponse structs.MapboxResponse
	if err := json.Unmarshal(body, &mapboxResponse); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return nil, err
	}
	if len(mapboxResponse.Features) == 0 {
		return nil, fmt.Errorf("no features found in response for organisation %s", organisation.Name)
	}

	return &mapboxResponse, nil
}
