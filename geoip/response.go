package geoip

type Response struct {
	Title          string      `json:"title"`
	IsRussia       string      `json:"isRussia"`
	IsNewRussia    string      `json:"isNewRussia"`
	OsmAddress     string      `json:"osmAddress"`
	Coordinates    Coordinates `json:"coordinates"`
	CountryName    string      `json:"countryName"`
	CountryIsoCode string      `json:"countryIsoCode"`
}

type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}
