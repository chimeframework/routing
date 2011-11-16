package routing

import (
	// "gorilla.googlecode.com/hg/gorilla/mux"
)

type Route struct {
	// mux.Route
}

func NewRoute(name string) *Route {
	this := &Route{}
	return this
}

type RouteCollection struct {
    
}

func NewRouteCollection() *RouteCollection {
    return &RouteCollection{}
}

