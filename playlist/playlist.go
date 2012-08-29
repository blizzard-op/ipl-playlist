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
	p.Items = make([]*PlaylistBlock, count)

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

		p.Items[i] = new(PlaylistBlock).Init(title, series, filepathsNode.(yaml.List))
	}
 	return p
}

func (p *Playlist) TotalItems() int {
	total := 0
	for _, block := range p.Items {
		total += len(block.Items)
	}
	return total
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

func (p *Playlist) Make() (*os.File, error) {
	log.Printf("Making playlist...")
	log.Printf("total = %d", p.TotalItems())

	tracks := make([]XspfTrack, p.TotalItems())
	index := 0
	for _, block := range p.Items {
		for _, item := range block.Items {
			tracks[index] = XspfTrack{Location: "file://" + item.Name()}
			index++
		}
	}

	xspf := &XspfPlaylist{Version: "1", Xmlns: "http://xspf.org/ns/0/", XspfTracks: tracks}

	xmlstring, err := xml.MarshalIndent(xspf, "", "    ")
	if err != nil {
	    log.Fatalf("xml.MarshalIndent: %v", err)
	}

	// create file
	output, err := os.Create("out.xspf")
	if err != nil {
	    log.Fatalf("os.Create: %v", err)
	}

	// write file
	bytesWritten, err := output.Write( []byte(xml.Header + string(xmlstring)) )
	if err != nil {
	    log.Fatalf("output.Write: %v", err)
	}
	log.Printf("bytesWritten: %d", bytesWritten)

	// close file
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
	log.Printf("Initializing block for %s", t)
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
	b.Duration = b.GetDuration()

	return b
}

func (b *PlaylistBlock) GetDuration() int {
	total := 0

	output_filepath := "tmp.flv"
	cleanup(output_filepath)

	// Duration: 00:08:59.66
	exp, err := regexp.Compile("Duration: ([0-9]{2}:[0-9]{2}:[0-9]{2}).[0-9]{2}")
	if err != nil {
		log.Fatalf("regexp.Compile: %v", err)
	}

	for _, f := range b.Items {
		path := f.Name()
		log.Printf("Getting duration for %s", path)

		cmd := exec.Command("ffmpeg", "-i", path, "-c", "copy", "-t", "1", output_filepath) // hack to get zero exit code

		stdout, er := cmd.CombinedOutput()
		if er != nil {
			log.Fatalf("cmd.CombinedOutput: %v", er)
		}

		//log.Printf("output: %s", stdout)

		result := exp.FindSubmatch(stdout)
		if result == nil {
			log.Fatalf("Could not determine duration")
		}
		log.Printf("%s", result)
		parts := strings.Split(string(result[1]), ":")
		for i, part := range parts {
			x, _ := strconv.Atoi(part)
			val := int( math.Pow(60, float64(len(parts[i+1:]))) ) * x
			total = total + val
		}

		// cleanup
		cleanup(output_filepath)
	}

	log.Printf("total: %d", total)
	return total
}

func cleanup(path string) {
	if f, err := os.Open(path); err == nil {
		if err := os.Remove(f.Name()); err != nil {
			log.Fatal(err)
		}
	}
}