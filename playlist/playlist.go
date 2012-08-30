package playlist

import (
	"fmt"
	"time"
	"log"
	"os"
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

func (p *Playlist) Output() (*os.File, error) {
	fmt.Println("Outputting playlist...")
	items := p.ArrangedItems()
	tracks := make([]XspfTrack, 0)
	for _, block := range items {
		for _, item := range block.Items {
			tracks = append(tracks, XspfTrack{Location: "file://" + item.Name()})
		}
	}
	x := XspfPlaylist{Version: "1", Xmlns: "http://xspf.org/ns/0/", XspfTracks: tracks}
	return x.Output()
}
