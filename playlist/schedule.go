package playlist

import (
	"fmt"
	"math"
	"time"
)

func (p *Playlist) availableDuration() time.Duration {
	return p.EndsAt.Sub(p.StartsAt)
}

func (p *Playlist) ArrangedItems() []*PlaylistBlock {	
	fmt.Printf("Target duration: %s\n", p.availableDuration().String())
	d := int(p.availableDuration().Seconds())
	var items []*PlaylistBlock
	var block *PlaylistBlock
	var i, primaryIndex, extrasIndex int
	for {		
		if d < 60 {
			break // less than 1 minute of available duration left
		}

		if (math.Mod(float64(i), float64(2)) == 0) {
			block = p.nextBlockToFill(p.Items, primaryIndex, d)
			primaryIndex += 1
			if (primaryIndex >= len(p.Items)) {
				primaryIndex = 0
			}
		} else {
			block = p.nextBlockToFill(p.ExtraItems, extrasIndex, d)
			extrasIndex += 1
			if (extrasIndex >= len(p.ExtraItems)) {
				extrasIndex = 0
			}
		}

		if ( block != nil ) {
			fmt.Printf("Available=%ds; Arranging %s [%ds]\n", d, block.Title, block.Duration)
			items = append(items, block)
			d -= block.Duration
		}
		
		i += 1
	}
	fmt.Printf("Total items arranged: %d\n%ds out of %ds remaining\n", len(items), d, int(p.availableDuration()))
	return items
}

func (p *Playlist) nextBlockToFill(items []*PlaylistBlock, startingIndex int, duration int) *PlaylistBlock {
	i := startingIndex
	for {
		block := items[i]
		if block.Duration <= duration {
			return block
		}
		i += 1
		if i >= len(items) {
			i = 0
		}
		if i == startingIndex {
			break
		}
	}
	return nil
}