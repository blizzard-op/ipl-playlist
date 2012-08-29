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
	Items := make([]*PlaylistBlock, count)

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

		Items[i] = new(PlaylistBlock).Init(title, series, filepathsNode.(yaml.List))
	}
 	return p
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
		durationParts := strings.Split(string(result[1]), ":")
		log.Printf("%q", durationParts)

		// cleanup
		cleanup(output_filepath)
	}

	return 0
}

func cleanup(path string) {
	if f, err := os.Open(path); err == nil {
		if err := os.Remove(f.Name()); err != nil {
			log.Fatal(err)
		}
	}
}