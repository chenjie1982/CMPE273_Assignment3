package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"gopkg.in/mgo.v2/bson"
)

func MongoCreate(Information *Info) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer sess.Close()
	//fmt.Printf("open DB successfully\n")

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("cmpe273").C("Location")

	count,err := db.Find(nil).Count()
	Information.Id = count+1
	//Information.Id = bson.NewObjectId()
	//fmt.Println("Information.Id:", Information.Id);

	err = db.Insert(&Information)
	if err != nil {
		fmt.Printf("Can't insert DB: %v\n", err)
		os.Exit(1)
	}

	var results []Info
	err = db.Find(bson.M{"_id": Information.Id}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	//fmt.Println("Results All: ", results)

	err = db.Find(bson.M{}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	//fmt.Println("Results All: ", results)
}

func MongoQuery(Id int/*bson.ObjectId */) (Info, error) {

    var result Info
	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return result, err
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")


	err = collection.Find(bson.M{"_id": Id}).One(&result)

	if err != nil {
		fmt.Println("Location ID Not find: ", Id)
		return result, err
	}
	//fmt.Println("Results : ", result)

	return result,nil
}

func MongoUpdate(Information Info ) (Info, error) {

	var LocalInfo Info

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return LocalInfo,err;
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")

	colQuerier := bson.M{"_id": Information.Id}

	change := bson.M{"$set": bson.M{"address": Information.Address,
									"city": Information.City,
									"state": Information.State,
									"zip": Information.Zip,
									"coordinate": bson.M{"lat":Information.Coordinate.Lat,
											"lng":Information.Coordinate.Lng}}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
		return LocalInfo,err
	}
	Info,error := MongoQuery(Information.Id)
	return Info,error
}




func MongoRemove(Id int/*bson.ObjectId*/ ) (error) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error ", err)
		panic(err)
		return err;
	}
	defer sess.Close()
	//fmt.Printf("MongoRemove: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")

	err = collection.Remove(bson.M{"_id": Id})
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func MongoQueryAllLocations() (AllInfo, error) {

	var result AllInfo
	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return result, err
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")


	err = collection.Find(nil).All(&result.AllLocations)

	if err != nil {
		panic(err)
		return result, err
	}

	return result,nil
}

func MongoCreateTrip(Information *TripInfo) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer sess.Close()
	//fmt.Printf("open DB successfully\n")

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("cmpe273").C("Trip")

	count,err := db.Find(nil).Count()
	Information.ID = count+1
	//Information.Id = bson.NewObjectId()
	//fmt.Println("Information.Id:", Information.Id);

	err = db.Insert(&Information)
	if err != nil {
		fmt.Printf("Can't insert DB: %v\n", err)
		os.Exit(1)
	}
}

func MongoUpdateTripState(status string, Id int ) ( error) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return err;
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Trip")

	colQuerier := bson.M{"_id": Id}

	change := bson.M{"$set": bson.M{"status": status}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func MongoQueryTrip(Id int/*bson.ObjectId */) (TripInfo, error) {

	var result TripInfo
	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return result, err
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Trip")


	err = collection.Find(bson.M{"_id": Id}).One(&result)

	if err != nil {
		//panic(err)
		return result, err
	}
	fmt.Println("Results : ", result)

	return result,nil
}

func MongoTripRequestFind(Id int) (TripRequest, error){

	var requesttrip TripRequest

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer sess.Close()
	//fmt.Printf("open DB successfully\n")

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("TripRequest")

	err = collection.Find(bson.M{"_id": Id}).One(&requesttrip)
	if err != nil {
		fmt.Printf("Not Find Request: ", requesttrip)
		return requesttrip, err
	}
	return requesttrip,nil
}

func MongoTripRequestCreate(Information TripRequest) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer sess.Close()
	fmt.Printf("open DB successfully\n")

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("cmpe273").C("TripRequest")

	fmt.Println("Information:", Information);

	err = db.Insert(&Information)
	if err != nil {
		fmt.Printf("Can't insert DB: %v\n", err)
		os.Exit(1)
	}
}

func MongoTripRequestUpdate(Information TripRequest) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer sess.Close()
	fmt.Printf("open DB successfully\n")

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("cmpe273").C("TripRequest")

	fmt.Println("Information:", Information);

	colQuerier := bson.M{"_id": Information.ID}

	change := bson.M{"$set": bson.M{"status": Information.Status,
		"next_destination_location_id": Information.NextLocaId,
		"uber_wait_time_eta": Information.WaitTime}}

	err = db.Update(colQuerier, change)
	if err != nil {
		panic(err)
		return
	}
	return

}