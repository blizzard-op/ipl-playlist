package main

import (
	"flag"
	"time"
	"os/exec"
	"log"
	playlist "github.com/ign/ipl-playlist/playlist"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

const timeFormat = time.RFC3339

var calendarName string; // name of Google Calendar channel

var startsAt time.Time; // time to start
var startsAtTime string;

var endsAt time.Time; // time to end
var endsAtTime string;

var configFilepath string;
var extrasConfigFilepath string;

func init() {
	now := time.Now()
	flag.StringVar(&calendarName, "calendar", "ignproleague_dev", "Name of channel. Default is ignproleague_dev.")
	flag.StringVar(&startsAtTime, "start", now.Format(timeFormat), "Start time. Default is now.")
	flag.StringVar(&endsAtTime, "end", now.Add(time.Hour*24).Format(timeFormat), "End time. Default is 24 hours from now.")
	flag.StringVar(&configFilepath, "config", "config.yml", "Config filepath. Default is './config.yml.'")
	flag.StringVar(&extrasConfigFilepath, "extras", "config.yml", "Extras config filepath. Default is './config.yml.'")
	flag.Parse() // parses the flags
	parseTimeVar(timeFormat, startsAtTime, &startsAt) // parse startsAt
	parseTimeVar(timeFormat, endsAtTime, &endsAt) // parse endsAt
}

func parseTimeVar(format string, value string, ptr *time.Time) {
	t, err := time.Parse(format, value)
    if err != nil {
        log.Fatalf("time.Parse: %v", err)
    }
    *ptr = t
}

func main() {
	log.Println("Starting...")

	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatalf("Could not find ffmpeg.")
	}

	config := yaml.ConfigFile(configFilepath)
	extrasConfig := yaml.ConfigFile(extrasConfigFilepath)

	p := new(playlist.Playlist).Init(startsAt, endsAt, *config, *extrasConfig)
	items := p.ArrangedItems()
	outfile, err := playlist.OutputPlaylist(items)
	if err != nil {
		log.Fatalf("Could not make playlist. %v", err)
	}
	log.Printf("Outputted playlist to %s", outfile.Name())
	p.Publish(calendarName, items)

	log.Println("Done.")
}