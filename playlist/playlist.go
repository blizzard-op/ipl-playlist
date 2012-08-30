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
	p.setItems( &p.Items, p.Config, "items" )
	p.setItems( &p.ExtraItems, p.ExtrasConfig, "extras" )
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

func (p *Playlist) setItems( items *[]*PlaylistBlock, config yaml.File, key string ) {
	lst, err := validateItems(config, key)
	if( err != nil){
		log.Fatalf("Failed items validation. %v", err)
	}
	*items = make([]*PlaylistBlock, lst.Len())
	for i, e := range lst {
		itemKey := key + "[" + strconv.Itoa(i) + "]"
		title, err := config.Get(itemKey + ".title")
		if (err != nil) {
			log.Fatalf("Missing title.")
		}
		series, err := p.Config.Get(itemKey + ".series")
		if (err != nil) {
			series = ""
		}
		filepathsNode, err := yaml.Child( e, "filepaths" )
		if err != nil {
			log.Fatalf("Missing filepaths for %s.", title)
		}
		(*items)[i] = new(PlaylistBlock).Init(title, series, filepathsNode.(yaml.List))
	}
}