package playlist

import (
	"fmt"
	"time"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

type Playlist struct {
	StartsAt, EndsAt time.Time
	Config yaml.File
}

func (p *Playlist) Make() {
	fmt.Println("Making playlist...")

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

	// get items list
	lst, ok := node.(yaml.List)
	if !ok {
		panic("Not a valid list of items.")
	}

	for _, item := range lst {
		fmt.Printf("- %v\n", item)
	}
}