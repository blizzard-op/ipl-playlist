package playlist

import (
	"strconv"
	"errors"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

func validateItems(config yaml.File, key string) (yaml.List, error) {
	node, err := yaml.Child( config.Root, key )
	if err != nil {
		return nil, err
	}
	lst, ok := node.(yaml.List)
	if (!ok || (lst.Len() <= 0)) {
		return nil, errors.New("Invalid items")
	}
	return lst, nil
}

func getItems(config yaml.File, key string) ([]*AvailableBlock, error){
	lst, err := validateItems(config, key)
	if( err != nil){
		return nil, errors.New("Invalid items")
	}
	items := make([]*AvailableBlock, lst.Len())
	for i, e := range lst {
		itemKey := key + "[" + strconv.Itoa(i) + "]"
		title, err := config.Get(itemKey + ".title")
		if (err != nil) {
			return nil, errors.New("Missing title")
		}
		series, err := config.Get(itemKey + ".series")
		if (err != nil) {
			series = ""
		}
		filepathsNode, err := yaml.Child( e, "filepaths" )
		if err != nil {
			return nil, errors.New("Missing filepaths for " + title)
		}
		publish := true
		if (key == "extras") {
			publish = false
		}
		items[i] = new(AvailableBlock).Init(title, series, filepathsNode.(yaml.List), publish)
	}
	return items, nil
}