package playlist

import (
	"fmt"
	"time"
	"strconv"
	"os"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

// Playlist is a description of a set of media files.
type Playlist struct {
	StartsAt, EndsAt time.Time
	Config yaml.File
	Items []*PlaylistBlock
}

func (p *Playlist) Init(s time.Time, e time.Time, c yaml.File) *Playlist {
 	fmt.Println("Initializing playlist...")
 	p.StartsAt = s
 	p.EndsAt = e
 	p.Config = c

	// items	
	node, err := yaml.Child( p.Config.Root, "items" )
	if err != nil {
		panic(err) // items node must be present
	}
	lst, ok := node.(yaml.List)
	if !ok {
		panic("Not a valid list of items.")
	}
	count := lst.Len()
	if (count <= 0) {
		panic("List of items is empty.") // items node must be a non-empty list
	}
	Items := make([]*PlaylistBlock, count)

	// blocks
	for i, e := range lst {
		itemKey := "items[" + strconv.Itoa(i) + "]"

		title, err := p.Config.Get(itemKey + ".title")
		if (err != nil) {
			panic(err)
		}

		series, err := p.Config.Get(itemKey + ".series")
		if (err != nil) {
			panic(err)
		}

		filepathsNode, err := yaml.Child( e, "filepaths" )
		if err != nil {
			panic(err)
		}

		Items[i] = new(PlaylistBlock).Init(title, series, filepathsNode.(yaml.List))
	}
 	return p
}

// A PlaylistBlock is a description of a set of related media files to keep grouped together.
type PlaylistBlock struct {
	Title string
	Series string
	Items []*os.File
}

func (b *PlaylistBlock) Init(t string, s string, filepaths yaml.List) *PlaylistBlock {
	b.Title = t
	b.Series = s
	count := filepaths.Len()
	if (count <= 0) {
		panic("List of filepaths is empty.")
	}
	b.Items = make([]*os.File, count)

	for i, e := range filepaths {
		f, err := os.Open(e.(yaml.Scalar).String())
		if err != nil {
			panic(err)
		}
		b.Items[i] = f
	}

	return b
}