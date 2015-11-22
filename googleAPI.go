package main
import (
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type LocationInfo struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
			       Lat float64 `json:"lat"`
			       Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport  struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

func QueryInfo(Address string) (Info, error) {
	//create Yahoo YQL to query database
	urlPath :=  "http://maps.google.com/maps/api/geocode/json?address="
	//convert space in url to UTF8
	urlPath += url.QueryEscape(Address)
	urlPath += "&sensor=false"

	fmt.Println(urlPath)
	var locationInfo Info
	var l LocationInfo
	//sent the http request to Google server
	res, err := http.Get(urlPath)
	if err!=nil {
		fmt.Println("QueryInfo: http.Get",err)
		return locationInfo,err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("QueryInfo: ioutil.ReadAll",err)
		return locationInfo,err
	}
	//parser the json that yahoo server responses, and store in the q

	err = json.Unmarshal(body, &l)

	if err!=nil {
		fmt.Println("QueryInfo: json.Unmarshal",err)
		return locationInfo,err
	}

	locationInfo.Coordinate.Lat = l.Results[0].Geometry.Location.Lat;
	locationInfo.Coordinate.Lng = l.Results[0].Geometry.Location.Lng;


	return locationInfo,nil

}
