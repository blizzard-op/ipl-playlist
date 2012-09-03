package playlist

import (
	"fmt"
	"time"
	"log"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

func (p *Playlist) Init(s time.Time, e time.Time, c yaml.File, xc yaml.File) *Playlist {
 	log.Println("Initializing playlist...")
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
