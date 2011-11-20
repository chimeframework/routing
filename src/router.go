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
}

func NewRouter() *Router {
	return &Router{}
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
