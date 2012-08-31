package playlist

import (
	"fmt"
	"reflect"
)

func (p *Playlist) Publish(items []*PlaylistBlock) () {
	fmt.Println("Publishing playlist...")
	for _, block := range items {
		fmt.Printf("Publishing %s\n", block.Title)
	}
	return
}