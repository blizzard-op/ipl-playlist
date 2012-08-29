package main

import (
	"fmt"
	"flag"
	"time"
	"os/exec"
	"log"
	playlist "./playlist"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

const timeFormat = time.RFC3339

// franchise
var franchiseName string;

// start time
var startsAt time.Time;
var startsAtTime string;

// end time
var endsAt time.Time;
var endsAtTime string;

// config filepath
var configFilepath string;

func init() {
	now := time.Now()
	flag.StringVar(&franchiseName, "franchise", "StarCraft 2", "Name of franchise. Default is StarCraft 2.")
	flag.StringVar(&startsAtTime, "start", now.Format(timeFormat), "Start time. Default is now.")
	flag.StringVar(&endsAtTime, "end", now.Add(time.Hour * 24).Format(timeFormat), "End time. Default is 24 hours from now.")
	flag.StringVar(&configFilepath, "config", "config.yml", "Config filepath. Default is './config.yml.'")
	flag.Parse() // parses the flags
	parseTimeVar(timeFormat, startsAtTime, &startsAt) // parse startsAt
	parseTimeVar(timeFormat, endsAtTime, &endsAt) // parse endsAt
}

func parseTimeVar(format string, value string, ptr *time.Time) {
	t, err := time.Parse(format, value)
    if err != nil {
        fmt.Println(err)
        return
    }
    *ptr = t
}

func main() {
	fmt.Println("Starting...")

	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("Could not find ffmpeg. Please")
	}

	config := yaml.ConfigFile(configFilepath)

	playlist := new(playlist.Playlist).Init(startsAt, endsAt, *config)
	output, err := playlist.Make()
	if err != nil {
		log.Fatalf("Could not make playlist. %v", err)
	}
	log.Printf("%T", output)

	fmt.Println("Done.")
}