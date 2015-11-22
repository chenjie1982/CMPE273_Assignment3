package main
import (
	"fmt"
	"math"
	"net/http"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"bytes"
)

type TripInfo struct {
	ID                     int     `json:"id"  bson:"_id,omitempty"`
	Status                 string  `json:"status" bson:"status"`
	StartingFromLocationID int     `json:"starting_from_location_id" bson:"starting_from_location_id"`
	BestRouteLocationIds   []int   `json:"best_route_location_ids" bson:"best_route_location_ids"`
	TotalUberCosts         int     `json:"total_uber_costs" bson:"total_uber_costs"`
	TotalUberDuration      int     `json:"total_uber_duration" bson:"total_uber_duration"`
	TotalDistance          float64 `json:"total_distance" bson:"total_distance"`
}

type TripRequest struct {
	ID                     int     `json:"id"  bson:"_id,omitempty"`
	Status                 string  `json:"status" bson:"status"`
	StartingFromLocationID int     `json:"starting_from_location_id" bson:"starting_from_location_id"`
	NextLocaId             int     `json:"next_destination_location_id" bson:"next_destination_location_id"`
	BestRouteLocationIds   []int   `json:"best_route_location_ids" bson:"best_route_location_ids"`
	TotalUberCosts         int     `json:"total_uber_costs" bson:"total_uber_costs"`
	TotalUberDuration      int     `json:"total_uber_duration" bson:"total_uber_duration"`
	TotalDistance          float64 `json:"total_distance" bson:"total_distance"`
	WaitTime               int     `json:"uber_wait_time_eta" bson:"uber_wait_time_eta"`
}

type UberPrice struct {
	Prices []struct {
		CurrencyCode    string  `json:"currency_code"`
		DisplayName     string  `json:"display_name"`
		Distance        float64 `json:"distance"`
		Duration        int     `json:"duration"`
		Estimate        string  `json:"estimate"`
		HighEstimate    int     `json:"high_estimate"`
		LowEstimate     int     `json:"low_estimate"`
		ProductID       string  `json:"product_id"`
		SurgeMultiplier int     `json:"surge_multiplier"`
	} `json:"prices"`
}

type TripPrice struct {
		Distance        float64 `json:"distance"`
		Duration        int     `json:"duration"`
		HighEstimate    int     `json:"high_estimate"`
		LowEstimate     int     `json:"low_estimate"`
}

type TiemEstimate struct {
	Driver          interface{} `json:"driver"`
	Eta             int         `json:"eta"`
	Location        interface{} `json:"location"`
	RequestID       string      `json:"request_id"`
	Status          string      `json:"status"`
	SurgeMultiplier int         `json:"surge_multiplier"`
	Vehicle         interface{} `json:"vehicle"`
}

type RideRequest struct {
	EndLatitude    float64 `json:"end_latitude"`
	EndLongitude   float64 `json:"end_longitude"`
	ProductID      string  `json:"product_id"`
	StartLatitude  float64 `json:"start_latitude"`
	StartLongitude float64 `json:"start_longitude"`
}


type UberProduct struct {
	Products []struct {
		Capacity     int    `json:"capacity"`
		Description  string `json:"description"`
		DisplayName  string `json:"display_name"`
		Image        string `json:"image"`
		PriceDetails struct {
			             Base            float64 `json:"base"`
			             CancellationFee int     `json:"cancellation_fee"`
			             CostPerDistance float64 `json:"cost_per_distance"`
			             CostPerMinute   float64 `json:"cost_per_minute"`
			             CurrencyCode    string  `json:"currency_code"`
			             DistanceUnit    string  `json:"distance_unit"`
			             Minimum         int     `json:"minimum"`
			             ServiceFees     []struct {
				             Fee  int    `json:"fee"`
				             Name string `json:"name"`
			             } `json:"service_fees"`
		             } `json:"price_details"`
		ProductID string `json:"product_id"`
	} `json:"products"`
}

func (result *TripPrice) QueryUberPrice(startLocationId int, endLocationId int) (error) {

	//var  TripPrice

	start,_ := MongoQuery(startLocationId)
	end,_ := MongoQuery(endLocationId)
	fmt.Println("start",start)
	fmt.Println("end",end)
	urlPath :=  "https://api.uber.com/v1/estimates/price?"
	parameters := "server_token="+"1Ga7UBb18RiYYQ_dug-EzsC1cYoG1xqyGH-kCBxv"

	parameters += "&start_latitude="+strconv.FormatFloat(start.Coordinate.Lat, 'f', 6, 64)
	parameters += "&start_longitude="+strconv.FormatFloat(start.Coordinate.Lng, 'f', 6, 64)
	parameters += "&end_latitude="+strconv.FormatFloat(end.Coordinate.Lat, 'f', 6, 64)
	parameters += "&end_longitude="+strconv.FormatFloat(end.Coordinate.Lng, 'f', 6, 64)
	urlPath += parameters

	fmt.Println(urlPath)

	var prices UberPrice

	//sent the http request to Google server
	res, err := http.Get(urlPath)
	if err!=nil {
		fmt.Println("QueryInfo: http.Get",err)
		return err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("QueryInfo: ioutil.ReadAll",err)
		return err
	}
	err = json.Unmarshal(body, &prices)
	//fmt.Println(prices)
	var cost = math.MaxInt64
	var id = 0

	if len(prices.Prices) == 0 {
		return nil
	}
	for i:= 0; i < len(prices.Prices); i++ {
		//fmt.Println("prices.Prices[i].LowEstimate: ",prices.Prices[i].LowEstimate)
		if (prices.Prices[i].LowEstimate != 0 && prices.Prices[i].LowEstimate < cost) {
			cost = prices.Prices[i].LowEstimate
			id = i;
			//productid = prices.Prices[i].ProductID
		}
	}
	result.Distance = prices.Prices[id].Distance
	result.Duration = prices.Prices[id].Duration
	result.LowEstimate = prices.Prices[id].LowEstimate
	result.HighEstimate = prices.Prices[id].HighEstimate
	//fmt.Println("result: ",result)
	return nil
}

func (result *TiemEstimate) RequestUber(startLocationId int, endLocationId int) error {

	//var result TiemEstimate

	start,_ := MongoQuery(startLocationId)
	end,_ := MongoQuery(endLocationId)

	//fmt.Println("start",start)
	//fmt.Println("end",end)

	urlPath :=  "https://sandbox-api.uber.com/v1/requests"

	var data RideRequest
	data.StartLatitude = start.Coordinate.Lat
	data.StartLongitude = start.Coordinate.Lng
	data.EndLatitude = end.Coordinate.Lat
	data.EndLongitude = end.Coordinate.Lng
	data.ProductID,_ = AvailableUberProductId(data);
	//fmt.Println("RequestUber: RideRequest",data)

	requestbody, _ := json.Marshal(data);
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(requestbody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiI3NzQ3OTY5ZS01YWZkLTQ3NjUtODI5Ni04NGUwZmNiZTZiMjMiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6ImJjZmYxZWM5LTRjZjgtNGU1Zi05MTVhLWI0ODhiYjBlNTQ3ZCIsImV4cCI6MTQ1MDU4NTI5OCwiaWF0IjoxNDQ3OTkzMjk4LCJ1YWN0IjoiTWNRN1h6dW9uVzdobmhTRWI2UkNDeFdHU2Q3eDl6IiwibmJmIjoxNDQ3OTkzMjA4LCJhdWQiOiJsNUxsVW5XX3lSbGFwVVpMdTd6ZDlIbjA0dVZHOUxMdyJ9.axozT_z8h4dbi_BbaUNEbUY4820J9a9uNJqjCEhXfkC-BAGHuBXWnSEtZx_czobME4PlvcqkbPvF1djig1Tlho31pN8HGIXuUIWLXN8Pzs5KYwKsmR7WymCf4sbHQ4yx6Qmd_3W7oidbsOerOLfIxpIrdzB2Z37gGPWjqfBy16K9TkzeBenouXPrQMJlfhkb6HXxSPCDiOj6F_QWhEj9lvmp_yDVmBAoP-8ByMy0PcW-MVZWb4k-2R20tKTPNLVrG27DCX5P18spH1xZVuRl4XVcwbapBbaVHLpOqSDxI4MPSwffA-z2rq16OHKT744OsOyRcB-d4WJ66auahCXslw"
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization","Bearer "+ token)
	//fmt.Println("RequestUber request:", req)
	res, err := client.Do(req)
	//fmt.Println("RequestUber response:", res)
	if err!=nil {
		fmt.Println("QueryInfo: http.Get",err)
		return err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("QueryInfo: ioutil.ReadAll",err)
		return err
	}
	err = json.Unmarshal(body, &result)
	//fmt.Println("RequestUber:", result.Eta)

	return nil
}

func AvailableUberProductId(data RideRequest) (string,error) {

	var result UberProduct


	fmt.Println("data",data)

	urlPath :=  "https://api.uber.com/v1/products?"
	parameters := "server_token="+"1Ga7UBb18RiYYQ_dug-EzsC1cYoG1xqyGH-kCBxv"

	parameters += "&latitude="+strconv.FormatFloat(data.StartLatitude, 'f', 6, 64)
	parameters += "&longitude="+strconv.FormatFloat(data.StartLongitude, 'f', 6, 64)

	urlPath += parameters

	//fmt.Println(urlPath)

	//sent the http request to Google server
	res, err := http.Get(urlPath)
	if err!=nil {
		fmt.Println("QueryInfo: http.Get",err)
		return "",err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("QueryInfo: ioutil.ReadAll",err)
		return "",err
	}
	err = json.Unmarshal(body, &result)
	//fmt.Println(prices)

	return result.Products[0].ProductID,nil
}
