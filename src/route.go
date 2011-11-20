package routing

import (
	"chime/components/config"
	"fmt"
	"gorilla.googlecode.com/hg/gorilla/mux"
	"regexp"
	"strings"
)

type Route struct {
	mux.Route
	pattern         string
	compiledPattern string
	defaults        map[string]interface{}
	requirements    map[string]interface{}
	options         map[string]interface{}
	hasCompiled     bool
}

func NewRoute(pattern string, defaults map[string]interface{}, requirements map[string]interface{}, options map[string]interface{}) *Route {
	this := &Route{}
	return this
}

func (this *Route) GetPattern() string {
	if !this.hasCompiled {
		this.Compile()
	}
	return this.compiledPattern
	// return this.pattern
}

// SetPattern sets the pattern for this route.
// This method implements a fluent interfafce
func (this *Route) SetPattern(pattern string) *Route {
	pattern = strings.TrimSpace(pattern)

	// patter must start with a slash '/'; if not add one
	if len(pattern) <= 0 || !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	this.pattern = pattern
	this.compiledPattern = ""
	this.hasCompiled = false
	return this
}

func (this *Route) SetOptions(options map[string]interface{}) *Route {
	this.options = options
	return this
}

func (this *Route) SetOption(name string, option interface{}) *Route {
	this.options[name] = option
	return this
}

func (this *Route) GetOption(name string) (interface{}, bool) {
	ret, ok := this.options[name]
	return ret, ok
}

func (this *Route) SetDefaults(defaults map[string]interface{}) *Route {
	this.defaults = defaults
	return this
}

func (this *Route) GetDefault(name string) (interface{}, bool) {
	ret, ok := this.defaults[name]
	return ret, ok
}
func (this *Route) SetDefault(name string, value interface{}) *Route {
	this.defaults[name] = value
	return this
}

func (this *Route) HasDefault(name string) bool {
	_, ok := this.GetDefault(name)
	return ok
}

func (this *Route) GetRequirements() map[string]interface{} {
	return this.requirements
}

func (this *Route) SetRequirements(requirements map[string]interface{}) *Route {
	this.requirements = make(map[string]interface{})
	for key, regex := range requirements {
		this.requirements[key] = sanitizeRequirements(key, regex)
	}
	return this
}

func (this *Route) GetRequirement(name string) (interface{}, bool) {
	ret, ok := this.requirements[name]
	return ret, ok
}
func (this *Route) SetRequirement(name string, value interface{}) *Route {
	this.requirements[name] = sanitizeRequirements(name, value)
	return this
}

func (this *Route) Compile() {
	if this.hasCompiled {
		return
	}
	reg := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	matches := reg.FindAllStringSubmatch(this.pattern, -1)
	this.compiledPattern = this.pattern
	this.hasCompiled = true
	if matches != nil {
		for _, match := range matches {
			placeholder := match[0]
			text := match[1]

			// check for requirement
			if req, ok := this.GetRequirement(text); ok {
				replacement := fmt.Sprintf("{%v:%v}", text, req)
				this.compiledPattern = strings.Replace(this.compiledPattern, placeholder, replacement, -1)
			}
			// TODO: check for defaults
		}
	}
}

func sanitizeRequirements(key string, val interface{}) string {
	// TODO check for an array requirements

	regex := config.ToString(val)
	if strings.HasPrefix(regex, "^") {
		regex = regex[1:]
	}

	if strings.HasSuffix(regex, "$") {
		regex = regex[0 : len(regex)-1]
	}
	return regex
}

/// Route Collection

type RouteCollection struct {
	Routes map[string]*Route
	prefix string
	Parent *RouteCollection
}

func NewRouteCollection() *RouteCollection {
	this := &RouteCollection{}
	this.Routes = make(map[string]*Route)
	this.Parent = nil
	return this
}

func (this *RouteCollection) AddPrefix(prefix string) {
	// a prefix must not end with a slash
	if strings.HasSuffix(prefix, "/"){
		return
	}

	// a preffix must start with a slash
	if !strings.HasPrefix(prefix, "/"){
		prefix = fmt.Sprintf("%v/", prefix)
	}
	this.prefix = fmt.Sprintf("%v%v", prefix, this.prefix)

	for _, route := range this.Routes{
		route.SetPattern(fmt.Sprintf("%v%v", prefix, route.GetPattern()))
	}
}

func (this *RouteCollection) GetPrefix() string{
    return this.prefix;
}

func (this *RouteCollection) Add(name string, route *Route) {
	// TODO: Check for name with invalid characters
    this.Routes[name] = route
}

func (this *RouteCollection) AddCollectionWithPrefix(collection *RouteCollection, prefix string) {
	collection.Parent = this
	collection.AddPrefix(prefix)

	for name, route :=  range collection.Routes{
		this.Routes[name] = route
	}
}
