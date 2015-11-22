package main
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
"math"
)

type Info struct {
	//Id          bson.ObjectId `json:"id"      bson:"_id,omitempty"`
	Id          int    `json:"id"      bson:"_id,omitempty"`
	Name        string `json:"name"    bson:"name"`
	Address     string `json:"address" bson:"address"`
	City        string `json:"city"    bson:"city"`
	State       string `json:"state"   bson:"state"`
	Zip         string `json:"zip"     bson:"zip"`
	Coordinate struct {
		Lat float64 `json:"lat" bson:"lat"`
		Lng float64 `json:"lng" bson:"lng"`
	} `json:"coordinate"        bson:"coordinate"`
}

type AllInfo struct {
	AllLocations []Info `json:"alllocations"`
}
type AddressInfo struct {
	Name        string `json:"name"    bson:"name"`
	Address     string `json:"address" bson:"address"`
	City        string `json:"city"    bson:"city"`
	State       string `json:"state"   bson:"state"`
	Zip         string `json:"zip"     bson:"zip"`
}

type LocationId struct {
	LocationIds            []int `json:"location_ids"`
	StartingFromLocationID int   `json:"starting_from_location_id"`
}

type Service struct{}


func Create(w http.ResponseWriter, r *http.Request) {

	var args AddressInfo
	var reply Info

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &args); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal Json from create emulator request.", body)
		return
	}
	addr := args.Address+","+args.City+","+args.State+","+args.Zip

	Information,err := QueryInfo(addr);

	reply.Address = args.Address;
	reply.City = args.City;
	reply.State = args.State;
	reply.Zip = args.Zip;
	reply.Name = args.Name;
	reply.Coordinate.Lat = Information.Coordinate.Lat
	reply.Coordinate.Lng = Information.Coordinate.Lng
	MongoCreate(&reply)
	//fmt.Println("Information.Id:"+reply.Id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
	return
}

func Update(w http.ResponseWriter, r *http.Request) {

	var args AddressInfo
	var reply Info

	vars := mux.Vars(r)
	//reply.Id = bson.ObjectIdHex(vars["location_id"])
	reply.Id,_ = strconv.Atoi(vars["location_id"])

	fmt.Println(reply.Id)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &args); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal Json from create emulator request.", body)
		return
	}

	addr := args.Address+","+args.City+","+args.State+","+args.Zip

	Information,err := QueryInfo(addr);

	//fmt.Println(Information)

	reply.Address = args.Address;
	reply.City = args.City;
	reply.State = args.State;
	reply.Zip = args.Zip;
	reply.Coordinate.Lat = Information.Coordinate.Lat
	reply.Coordinate.Lng = Information.Coordinate.Lng

	Information, err = MongoUpdate(reply)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	reply.Name = Information.Name

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}

}

func Query(w http.ResponseWriter, r *http.Request) {

	var reply Info
	var err error
	vars := mux.Vars(r)
	//reply.Id = bson.ObjectIdHex(vars["location_id"])
	reply.Id,_ = strconv.Atoi(vars["location_id"])

	reply, err = MongoQuery(reply.Id )
	if err != nil {

		fmt.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {

	var reply Info

	vars := mux.Vars(r)
	//reply.Id = bson.ObjectIdHex(vars["location_id"])
	reply.Id,_ = strconv.Atoi(vars["location_id"])
	err := MongoRemove(reply.Id )
	if err != nil {

		fmt.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return

}


func PlanTrip(w http.ResponseWriter, r *http.Request) {

	var args LocationId

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &args); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal Json from create emulator request.", body)
		return
	}
	fmt.Println("Carry Para:",args)

	startLocationId := args.StartingFromLocationID
	var cost = math.MaxInt64
	var id = 0;
	var count = len(args.LocationIds);
	var price TripPrice
	var trip TripInfo
	trip.Status = "planning"
	trip.StartingFromLocationID = args.StartingFromLocationID
	var priceTemp TripPrice
	for count > 0 {
		for i := 0; i < len(args.LocationIds); i++ {
			//call uber API to calculate the shortest route
			err1 := priceTemp.QueryUberPrice(startLocationId,args.LocationIds[i]);
			if err1 != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(422) // unprocessable entity
				fmt.Printf(err.Error())
				return;
			}
			fmt.Println("price.LowEstimate:",priceTemp.LowEstimate,"[id]",i)
			fmt.Println("cost:",cost)
			if(cost > priceTemp.LowEstimate) {
				cost = priceTemp.LowEstimate;
				price.Distance = priceTemp.Distance;
				price.Duration = priceTemp.Duration;
				price.LowEstimate = priceTemp.LowEstimate;
				id = i;
			}
		}
		fmt.Println("ID***",id)
		fmt.Println("args.LocationIds[id]",args.LocationIds[id])
		trip.TotalDistance += price.Distance
		trip.TotalUberCosts += price.LowEstimate
		trip.TotalUberDuration += price.Duration
		//fmt.Println(trip.TotalDistance,trip.TotalUberCosts,trip.TotalUberDuration)
		fmt.Println("*********************")
		fmt.Println(args.LocationIds)
		trip.BestRouteLocationIds = append(trip.BestRouteLocationIds, args.LocationIds[id])
		startLocationId = args.LocationIds[id]
		args.LocationIds = append(args.LocationIds[:id], args.LocationIds[id+1:]...)
		//fmt.Println(startLocationId)
		fmt.Println(args.LocationIds)
		fmt.Println("*********************")
		cost = math.MaxInt64
		count--;
	}
	MongoCreateTrip(&trip)
	fmt.Println(trip)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(trip); err != nil {
		panic(err)
	}

}

func CheckTrip(w http.ResponseWriter, r *http.Request) {

	var trip TripInfo
	var err error
	vars := mux.Vars(r)
	Id,_ := strconv.Atoi(vars["trip_id"])

	trip, err = MongoQueryTrip(Id)
	if err != nil {

		fmt.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(trip); err != nil {
		panic(err)
	}
	w.Header().Set("Access-Control-Allow-Origin","*")
}

func QueryAllLocations(w http.ResponseWriter, r *http.Request) {

	Location, err := MongoQueryAllLocations()
	if err != nil {

		fmt.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin","*")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(Location); err != nil {
		panic(err)
	}
}

func RequestTrip(w http.ResponseWriter, r *http.Request) {

	var estimateTime TiemEstimate
	vars := mux.Vars(r)

	Id,_ := strconv.Atoi(vars["trip_id"])

	fmt.Println(Id)

	requesttrip, err := MongoTripRequestFind(Id)
	fmt.Println("requesttrip: ", requesttrip,err)
	//fmt.Println("requesttrip.status: ", requesttrip.Status)
	if err != nil {
		trip, err := MongoQueryTrip(Id)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		requesttrip.ID = trip.ID
		requesttrip.BestRouteLocationIds = trip.BestRouteLocationIds
		requesttrip.StartingFromLocationID = trip.StartingFromLocationID
		requesttrip.TotalDistance = trip.TotalDistance
		requesttrip.TotalUberCosts = trip.TotalUberCosts
		requesttrip.TotalUberDuration = trip.TotalUberDuration
		requesttrip.NextLocaId = trip.BestRouteLocationIds[0]
		err = estimateTime.RequestUber(trip.StartingFromLocationID,requesttrip.NextLocaId)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		requesttrip.WaitTime = estimateTime.Eta;
		requesttrip.Status = estimateTime.Status

		MongoTripRequestCreate(requesttrip)
		MongoUpdateTripState(requesttrip.Status,requesttrip.ID )

	} else if requesttrip.Status != "finish" {
		var i int
		for i = 0; i < len(requesttrip.BestRouteLocationIds); i++ {
			if(requesttrip.BestRouteLocationIds[i] == requesttrip.NextLocaId){
				break;
			}
		}
		fmt.Println("requesttrip.Status != finish: ",requesttrip)
		if( i+1 == len(requesttrip.BestRouteLocationIds)){
			requesttrip.Status = "finish"
			requesttrip.WaitTime = 0
		} else {
			requesttrip.NextLocaId = requesttrip.BestRouteLocationIds[i+1]
			err = estimateTime.RequestUber(requesttrip.BestRouteLocationIds[i],requesttrip.NextLocaId)
			fmt.Printf("RequestTrip[estimateTime]",estimateTime.Eta)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}

			requesttrip.WaitTime = estimateTime.Eta
		}
		MongoTripRequestUpdate(requesttrip)
		MongoUpdateTripState(requesttrip.Status,requesttrip.ID )

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(requesttrip); err != nil {
		panic(err)
	}
}
