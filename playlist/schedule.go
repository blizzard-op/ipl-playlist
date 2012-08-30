package playlist

import (
	"fmt"
	"os"
	"log"
	"math"
	"encoding/xml"
)

func (p *Playlist) availableDuration() float64 {
	return p.EndsAt.Sub(p.StartsAt).Seconds()
}

func (p *Playlist) ArrangedItems() []*PlaylistBlock {	
	fmt.Printf("Target duration: %s\n", p.EndsAt.Sub(p.StartsAt).String())

	d := int(p.availableDuration())

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

func (p *Playlist) Make() (*os.File, error) {
	fmt.Printf("Making playlist...\n")

	items := p.ArrangedItems()
	tracks := make([]XspfTrack, 0)

	for _, block := range items {
		for _, item := range block.Items {
			tracks = append(tracks, XspfTrack{Location: "file://" + item.Name()})
		}
	}

	xspfPlaylist := &XspfPlaylist{Version: "1", Xmlns: "http://xspf.org/ns/0/", XspfTracks: tracks}

	xmlstring, err := xml.MarshalIndent(xspfPlaylist, "", "    ")
	if err != nil {
	    log.Fatalf("xml.MarshalIndent: %v", err)
	}

	// create
	output, err := os.Create("out.xspf")
	if err != nil {
	    log.Fatalf("os.Create: %v", err)
	}

	// write
	_, err = output.Write( []byte(xml.Header + string(xmlstring)) )
	if err != nil {
	    log.Fatalf("output.Write: %v", err)
	}

	// close
	err = output.Close()
	if err != nil {
	    log.Fatalf("output.Close: %v", err)
	}

	return output, err
}