package playlist

import (
	"fmt"
)

func (p *Playlist) Publish(calendarName string, items []*PlaylistBlock) () {
	fmt.Println("Publishing playlist to ")
	for _, block := range items {
		if( block.DoPublish ){			
			block.Publish(calendarName)
		}
	}
	return
}

func (block *PlaylistBlock) Publish(calendarName string) {
	fmt.Printf("Publishing %s to %s\n", block.Title, calendarName)
}