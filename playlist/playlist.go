package playlist

import (
	"fmt"
	"time"
	"strconv"
	"log"
	"errors"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

// Playlist is a description of a set of media files.
type Playlist struct {
	StartsAt, EndsAt time.Time
	Config yaml.File
	Items []*PlaylistBlock
	ExtrasConfig yaml.File
	ExtraItems []*PlaylistBlock
}

func (p *Playlist) Init(s time.Time, e time.Time, c yaml.File, xc yaml.File) *Playlist {
 	fmt.Println("Initializing playlist...")
 	p.StartsAt = s
 	p.EndsAt = e
 	p.Config = c
 	p.ExtrasConfig = xc
 	var err error
	p.Items, err = getItems(p.Config, "items")
	if( err != nil){
		log.Fatalf("Invalid items. %v", err)
	}
	p.ExtraItems, err = getItems(p.ExtrasConfig, "extras")
	if( err != nil){
		log.Fatalf("Invalid extras. %v", err)
	}
 	return p
}

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

func getItems(config yaml.File, key string) ([]*PlaylistBlock, error){
	lst, err := validateItems(config, key)
	if( err != nil){
		return nil, errors.New("Invalid items")
	}
	items := make([]*PlaylistBlock, lst.Len())
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
		items[i] = new(PlaylistBlock).Init(title, series, filepathsNode.(yaml.List))
	}
	return items, nil
}