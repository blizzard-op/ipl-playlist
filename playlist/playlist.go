package playlist

import (
	"fmt"
	"time"
	"strconv"
	"log"
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

func (p *Playlist) setItems( items *[]*PlaylistBlock, config yaml.File, key string ) {

	node, err := yaml.Child( config.Root, key )
	if err != nil {
		log.Fatalf("No items. %v", err) // items node must be present
	}
	lst, ok := node.(yaml.List)
	if !ok {
		log.Fatalf("Invalid items. %v", err)
	}
	count := lst.Len()
	if (count <= 0) {
		log.Fatalf("No items. %v", err) // items node must be a non-empty list
	}
	*items = make([]*PlaylistBlock, count)

	// blocks
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