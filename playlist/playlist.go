package playlist

import (
	"fmt"
	"time"
	"strconv"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

type Playlist struct {
	StartsAt, EndsAt time.Time
	Config yaml.File
	items []PlaylistBlock
}

func (p *Playlist) Init(s time.Time, e time.Time, c yaml.File) *Playlist {
 	fmt.Println("Initializing playlist...")

 	p.StartsAt = s
 	p.EndsAt = e
 	p.Config = c

	// items node must be present
	node, err := yaml.Child( p.Config.Root, "items" )
	if err != nil {
		panic(err)
	}

	// items node must be a non-empty list
	itemsCount, err := p.Config.Count("items")
	if (err != nil) {
		panic(err)
	}
	if (itemsCount <= 0) {
		panic("List of items is empty.")
	}

	// construct items
	items := make([]PlaylistBlock, itemsCount)

	// get items list
	lst, ok := node.(yaml.List)
	if !ok {
		panic("Not a valid list of items.")
	}

	// blocks
	for i, _ := range lst {
		title, err := p.Config.Get("items[" + strconv.Itoa(i) + "].title")
		if (err != nil) {
			panic(err)
		}
		series, err := p.Config.Get("items[" + strconv.Itoa(i) + "].series")
		if (err != nil) {
			panic(err)
		}
		items[i] = PlaylistBlock{ title, series }
	}
 	return p
}

type PlaylistBlock struct {
	Title string
	Series string
}