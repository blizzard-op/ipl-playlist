package playlist

import (
	"fmt"
)

func (p *Playlist) Publish(items []*PlaylistBlock) () {
	fmt.Println("Publishing playlist...")
	for _, block := range items {
		if( block.DoPublish ){			
			block.Publish()
		}
	}
	return
}

func (block *PlaylistBlock) Publish() {
	fmt.Printf("Publishing %s...\n", block.Title)
}