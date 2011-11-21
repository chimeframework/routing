package routing

import (
	"chime/components/httpcontext"
	"gorilla.googlecode.com/hg/gorilla/context"
	"gorilla.googlecode.com/hg/gorilla/mux"
	"net/url"
)

type Router struct {
	*mux.Router
	// routeMaps map[string]*Route
	collection *RouteCollection
	loader *YamlFileLoader
	resource string
}

func NewRouter(loader *YamlFileLoader, resource string) *Router {
	return &Router{loader:loader, resource: resource}
}

func (this *Router) GenerateUrl(name string, pairs []string) *url.URL {
	return this.NamedRoutes[name].URL(pairs...)
}

func (this *Router) MatchRequest(request *httpcontext.Request) mux.RouteVars {
	rawRequest := request.Request
	if match, ok := this.Match(rawRequest); ok {
		defer context.DefaultContext.Clear(rawRequest)
		params := mux.Vars(rawRequest)
		params[httpcontext.ROUTE_PARAM] = match.Route.GetName()
		return params
	}

	panic("No match found")
}

func (this *Router) GetRouteCollection() *RouteCollection{
    if this.collection == nil{
    	this.collection = this.loader.Load(this.resource)
    }

    // TODO: Compile each route
    return this.collection
}

func (this *Router) Compile() {
    for name, route := range this.GetRouteCollection().Routes{
    	muxRoute := this.NewRoute().Name(name).Path(route.GetPattern())

    	if methods, ok := route.GetRequirement(ROUTE_REQUIREMENTS_METHOD); ok{
	    	muxRoute.Methods(methods.([]string)...)
    	}

    	if schemes, ok := route.GetRequirement(ROUTE_REQUIREMENTS_SCHEME); ok{
    		muxRoute.Schemes(schemes.([]string)...)
    	}
    }
}
