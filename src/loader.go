package routing

import (
	"chime/components/config"
	"chime/components/yaml"
	"fmt"
	"path"
	"sort"
	"strings"
)

const (
	ROUTE_RESOURCE     = "resource"
	ROUTE_TYPE         = "type"
	ROUTE_PREFIX       = "prefix"
	ROUTE_DEFAULTS     = "defaults"
	ROUTE_REQUIREMENTS = "requirements"
	ROUTE_OPTIONS      = "options"
	ROUTE_PATTERN      = "pattern"
)

type YamlFileLoader struct {
	locator        *config.FileLocator
	availableKeys  []string
	loadedResource map[string]bool
}

func NewYamlFileLoader(locator *config.FileLocator) *YamlFileLoader {
	this := &YamlFileLoader{locator: locator}
	// Note: the keys should be in ascending order to be able to search using sort.SearchStrings
	this.availableKeys = []string{ROUTE_DEFAULTS, ROUTE_OPTIONS, ROUTE_PATTERN, ROUTE_PREFIX, ROUTE_REQUIREMENTS, ROUTE_RESOURCE, ROUTE_TYPE}
	// This will hold the loaded resource so that we can avoid circular reference
	this.loadedResource = make(map[string]bool)
	return this
}

func (this *YamlFileLoader) Load(file string) *RouteCollection {
	fullPaths := this.locator.LocateFirst(file)
	fullPath := fullPaths[0]
	routeConfigs := yaml.Parse(fullPath)
	collection := NewRouteCollection()
	//TODO: Add resource to route collection

	for routeName, value := range routeConfigs {
		routeConfig := value.(map[interface{}]interface{})
		routeConfig = this.normalizeRouteConfig(routeConfig)

		if res, ok := routeConfig[ROUTE_RESOURCE]; ok {
			prefix, ok := routeConfig[ROUTE_PREFIX]
			if !ok {
				prefix = ""
			}

			currDir, _ := path.Split(fullPath)
			coll := this.importResource(config.ToString(res), file, currDir)
			collection.AddCollectionWithPrefix(coll, config.ToString(prefix))
			// TODO: Support for ROUTE TYPE
		} else {
			this.parseRoute(collection, config.ToString(routeName), routeConfig)
		}
	}
	return collection
}

func (this *YamlFileLoader) parseRoute(collection *RouteCollection, routeName string, routeConfigs map[interface{}]interface{}) {

	pattern, ok := routeConfigs[ROUTE_PATTERN]

	// must have 'ROUTE_PATTERN' section
	if !ok {
		panic(fmt.Sprintf("You must define a pattern for the %v route.", routeName))
	}

	tempMap, ok := routeConfigs[ROUTE_DEFAULTS]
	if !ok {
		tempMap = make(map[string]interface{})
	}
	defaults := tempMap.(map[string]interface{})

	tempMap, ok = routeConfigs[ROUTE_REQUIREMENTS]
	if !ok {
		tempMap = make(map[string]interface{})
	}
	requirements := tempMap.(map[string]interface{})

	tempMap, ok = routeConfigs[ROUTE_OPTIONS]
	if ok {
		tempMap = make(map[string]interface{})
	}

	options := tempMap.(map[string]interface{})

	route := NewRoute(config.ToString(pattern), defaults, requirements, options)
	collection.Add(routeName, route)
}

func (this *YamlFileLoader) importResource(resource string, sourceResource string, currDir string) *RouteCollection {
	fullPaths := this.locator.LocateFirstFrom(resource, currDir)

	resource = fullPaths[0]
	if _, ok := this.loadedResource[resource]; ok {
		panic(fmt.Sprintf("Circular reference detected for %v", resource))
	}

	this.loadedResource[resource] = true
	// this is a recursive call, but we have avoided a circular reference
	ret := this.Load(resource)
	this.loadedResource[resource] = false
	return ret
}

func (this *YamlFileLoader) normalizeRouteConfig(config map[interface{}]interface{}) map[interface{}]interface{} {
	for k, _ := range config {
		index := sort.SearchStrings(this.availableKeys, k.(string))

		if index >= len(this.availableKeys) || this.availableKeys[index] != k {
			panic(fmt.Sprintf("Yaml routing loader does not support given key: %v. Expected one of the (%v).", k, strings.Join(this.availableKeys, ", ")))
		}
	}

	return config
}
