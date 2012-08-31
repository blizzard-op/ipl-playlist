package playlist

import (
	"fmt"
)

func (p *Playlist) Publish(items []*PlaylistBlock) () {
	fmt.Println("Publishing playlist...")
	for _, block := range items {
		fmt.Printf("Publishing %s %v\n", block.Title, block.Publish)
	}
	return
}