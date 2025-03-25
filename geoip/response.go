// geoip/response.go
package geoip

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Title          *string      `json:"title"`
	IsRussia       *bool        `json:"isRussia"`
	IsNewRussia    *bool        `json:"isNewRussia"`
	OsmAddress     *string      `json:"osmAddress"`
	Coordinates    *Coordinates `json:"coordinates"`
	CountryName    string       `json:"country_name"`     // <-- updated
	CountryIsoCode string       `json:"country_iso_code"` // <-- updated
}

type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

func (c *Coordinates) UnmarshalJSON(data []byte) error {

	// First try array of floats: [longitude, latitude]
	var arr [2]float64
	if err := json.Unmarshal(data, &arr); err == nil {
		c.Long = fmt.Sprintf("%f", arr[0])
		c.Lat = fmt.Sprintf("%f", arr[1])
		return nil
	}

	// Then try object with named fields
	var obj struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		c.Lat = obj.Lat
		c.Long = obj.Long
		return nil
	}

	return fmt.Errorf("coordinates: unknown format: %s", string(data))
}
