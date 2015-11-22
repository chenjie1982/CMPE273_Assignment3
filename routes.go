package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Create",
		"POST",
		"/locations",
		Create,
	},
	Route{
		"Query",
		"GET",
		"/locations/{location_id}",
		Query,
	},
	Route{
		"Update",
		"PUT",
		"/locations/{location_id}",
		Update,
	},
	Route{
		"Remove",
		"DELETE",
		"/locations/{location_id}",
		Remove,
	},
	Route{
		"CheckTrip",
		"GET",
		"/trips/{trip_id}",
		CheckTrip,
	},
	Route{
		"PlanTrip",
		"POST",
		"/trips",
		PlanTrip,
	},
	Route{
		"QueryAllLocations",
		"GET",
		"/findlocations",
		QueryAllLocations,
	},
	Route{
		"RequestTrip",
		"POST",
		"/trips/{trip_id}/request",
		RequestTrip,
	},
}
