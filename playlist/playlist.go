package playlist

import (
	"fmt"
	"time"
	"strconv"
	"os"
	"os/exec"
	"log"
	"regexp"
	"strings"
	"math"
	"encoding/xml"
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
	p.setItems( &p.Items, p.Config, "items" )
	p.setItems( &p.ExtraItems, p.ExtrasConfig, "extras" )
 	return p
}

func (p *Playlist) setItems( items *[]*PlaylistBlock, config yaml.File, key string ) {

	node, err := yaml.Child( config.Root, key )
	if err != nil {
		log.Fatalf("No items. %v", err) // items node must be present
	}
	lst, ok := node.(yaml.List)
	if !ok {
		log.Fatalf("Invalid items. %v", err)
	}
	count := lst.Len()
	if (count <= 0) {
		log.Fatalf("No items. %v", err) // items node must be a non-empty list
	}
	*items = make([]*PlaylistBlock, count)

	// blocks
	for i, e := range lst {
		itemKey := key + "[" + strconv.Itoa(i) + "]"

		title, err := config.Get(itemKey + ".title")
		if (err != nil) {
			log.Fatalf("Missing title.")
		}

		series, err := p.Config.Get(itemKey + ".series")
		if (err != nil) {
			series = ""
		}

		filepathsNode, err := yaml.Child( e, "filepaths" )
		if err != nil {
			log.Fatalf("Missing filepaths for %s.", title)
		}

		(*items)[i] = new(PlaylistBlock).Init(title, series, filepathsNode.(yaml.List))
	}

}

type XspfPlaylist struct {
    XMLName xml.Name `xml:"playlist"`
    Version string `xml:"version,attr"`
    Xmlns string `xml:"xmlns,attr"`
    XspfTracks []XspfTrack `xml:"trackList>track"`
}

type XspfTrack struct {
	Location string `xml:"location"`
}

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

	xspf := &XspfPlaylist{Version: "1", Xmlns: "http://xspf.org/ns/0/", XspfTracks: tracks}

	xmlstring, err := xml.MarshalIndent(xspf, "", "    ")
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

// A PlaylistBlock is a description of a set of related media files to keep grouped together.
type PlaylistBlock struct {
	Title string
	Series string
	Items []*os.File
	Duration int
}

func (b *PlaylistBlock) Init(t string, s string, filepaths yaml.List) *PlaylistBlock {
	b.Title = t
	b.Series = s
	count := filepaths.Len()
	if (count <= 0) {
		log.Fatalf("No filepaths for %s.", t)
	}	
	b.Items = make([]*os.File, count)
	for i, e := range filepaths {
		f, err := os.Open(e.(yaml.Scalar).String())
		if err != nil {
			log.Fatalf("Missing file for %s.", t)
		}
		b.Items[i] = f
	}
	b.Duration = b.GetDuration()

	return b
}

func (b *PlaylistBlock) GetDuration() int {
	total := 0

	output_filepath := "tmp.flv" // TODO adapt to original file suffix
	cleanup(output_filepath)

	exp, err := regexp.Compile("Duration: ([0-9]{2}:[0-9]{2}:[0-9]{2}).[0-9]{2}")
	if err != nil {
		log.Fatalf("regexp.Compile: %v", err)
	}

	for _, f := range b.Items {
		path := f.Name()
		cmd := exec.Command("ffmpeg", "-i", path, "-c", "copy", "-t", "1", output_filepath) // hack to get zero exit code
		stdout, er := cmd.CombinedOutput()
		if er != nil {
			log.Fatalf("cmd.CombinedOutput: %v", er)
		}

		result := exp.FindSubmatch(stdout)
		if result == nil {
			log.Fatalf("Could not determine duration")
		}
		parts := strings.Split(string(result[1]), ":")
		for i, part := range parts {
			x, _ := strconv.Atoi(part)
			val := int( math.Pow(60, float64(len(parts[i+1:]))) ) * x
			total = total + val
		}

		cleanup(output_filepath)
	}
	fmt.Printf("Using %s [%ds]\n", b.Title, total)
	return total
}

func cleanup(path string) {
	if f, err := os.Open(path); err == nil {
		if err := os.Remove(f.Name()); err != nil {
			log.Fatal(err)
		}
	}
}