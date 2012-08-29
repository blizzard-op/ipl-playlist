package playlist

import (
	"fmt"
	"time"
	"strconv"
	"os"
	"os/exec"
	"log"
	//"io/ioutil"
	//"io"
	//"bufio"
	//"errors"
	//"bytes"
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

	for _, f := range b.Items {
		path := f.Name()
		log.Printf("Getting duration for %s", path)

		cmd := exec.Command("ffmpeg", "-i", path, "-c", "copy", "-t", "1", output_filepath) // hack to get zero exit code

		stdout, er := cmd.CombinedOutput()
		if er != nil {
			log.Fatalf("cmd.CombinedOutput: %v", er)
    	}

		// runError := cmd.Run()
		// if runError != nil {
		// 	log.Fatalf("cmd.Run: %v", runError)
  //   	}

		// b, readAllError := ioutil.ReadAll(stdout)
		// if readAllError != nil {
		// 	log.Fatalf("ioutil.ReadAll: %v", readAllError)
  //   	}
		log.Printf("output: %s", stdout)

		// cleanup
		cleanup(output_filepath)
	}

	return 0
}

func cleanup(path string) {
	if f, err := os.Open(path); err == nil {
		log.Printf("Cleaning up %s", f.Name())
		if err := os.Remove(path); err != nil {
			log.Fatal(err)
		}
	}
}