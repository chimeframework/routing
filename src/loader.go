package routing

import (
    "chime/components/config"
)

type YamlFileLoader struct {
    locator *config.FileLocator
}

func NewYamlFileLoader(locator *config.FileLocator) *YamlFileLoader {
    return &YamlFileLoader{locator:locator}
}

func (this *YamlFileLoader) Load(file string) *RouteCollection{
    collection := NewRouteCollection()
    return collection
}
