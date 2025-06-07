package structs

type MapboxResponse struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
}

type Feature struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	MapboxID       string      `json:"mapbox_id"`
	FeatureType    string      `json:"feature_type"`
	FullAddress    string      `json:"full_address"`
	Name           string      `json:"name"`
	NamePreferred  string      `json:"name_preferred"`
	Coordinates    Coordinates `json:"coordinates"`
	PlaceFormatted string      `json:"place_formatted"`
	MatchCode      MatchCode   `json:"match_code"`
	Context        Context     `json:"context"`
}

type Coordinates struct {
	Longitude      float64      `json:"longitude"`
	Latitude       float64      `json:"latitude"`
	Accuracy       string       `json:"accuracy"`
	RoutablePoints []RoutablePt `json:"routable_points"`
}

type RoutablePt struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MatchCode struct {
	AddressNumber string `json:"address_number"`
	Street        string `json:"street"`
	Postcode      string `json:"postcode"`
	Place         string `json:"place"`
	Region        string `json:"region"`
	Locality      string `json:"locality"`
	Country       string `json:"country"`
	Confidence    string `json:"confidence"`
}

type Context struct {
	Address      AddressComponent `json:"address"`
	Street       StreetComponent  `json:"street"`
	Neighborhood *Neighborhood    `json:"neighborhood,omitempty"`
	Postcode     NameOnly         `json:"postcode"`
	Place        NamedWithAlt     `json:"place"`
	District     Named            `json:"district"`
	Region       Region           `json:"region"`
	Country      Country          `json:"country"`
}

type AddressComponent struct {
	MapboxID      string `json:"mapbox_id"`
	AddressNumber string `json:"address_number"`
	StreetName    string `json:"street_name"`
	Name          string `json:"name"`
}

type StreetComponent struct {
	MapboxID string `json:"mapbox_id"`
	Name     string `json:"name"`
}

type Neighborhood struct {
	MapboxID   string    `json:"mapbox_id"`
	Name       string    `json:"name"`
	Alternate  *NameOnly `json:"alternate,omitempty"`
	WikidataID string    `json:"wikidata_id,omitempty"`
}

type NameOnly struct {
	MapboxID string `json:"mapbox_id"`
	Name     string `json:"name"`
}

type NamedWithAlt struct {
	MapboxID   string    `json:"mapbox_id"`
	Name       string    `json:"name"`
	WikidataID string    `json:"wikidata_id"`
	Alternate  *NameOnly `json:"alternate,omitempty"`
}

type Named struct {
	MapboxID   string `json:"mapbox_id"`
	Name       string `json:"name"`
	WikidataID string `json:"wikidata_id"`
}

type Region struct {
	MapboxID       string `json:"mapbox_id"`
	Name           string `json:"name"`
	WikidataID     string `json:"wikidata_id"`
	RegionCode     string `json:"region_code"`
	RegionCodeFull string `json:"region_code_full"`
}

type Country struct {
	MapboxID          string `json:"mapbox_id"`
	Name              string `json:"name"`
	WikidataID        string `json:"wikidata_id"`
	CountryCode       string `json:"country_code"`
	CountryCodeAlpha3 string `json:"country_code_alpha_3"`
}
