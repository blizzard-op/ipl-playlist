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

var skipPublish bool;
var skipOutput bool;

func init() {
	now := time.Now()
	flag.StringVar(&calendarName, "calendar", "ignproleague_dev", "Name of channel. Default is ignproleague_dev.")
	flag.StringVar(&startsAtTime, "start", now.Format(timeFormat), "Start time. Default is now.")
	flag.StringVar(&endsAtTime, "end", now.Add(time.Hour*2 + time.Minute*1).Format(timeFormat), "End time. Default is 24 hours from now.")
	flag.StringVar(&configFilepath, "config", "config.yml", "Config filepath. Default is './config.yml.'")
	flag.StringVar(&extrasConfigFilepath, "extras", "config.yml", "Extras config filepath. Default is './config.yml.'")
	flag.BoolVar(&skipPublish, "skipPublish", false, "Skip publishing. Default is false.")
	flag.BoolVar(&skipOutput, "skipOutput", false, "Skip output file. Default is false.")
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
		log.Fatal("Could not find ffmpeg.")
	}

	config := yaml.ConfigFile(configFilepath)
	extrasConfig := yaml.ConfigFile(extrasConfigFilepath)

	p := new(playlist.Playlist).Init(startsAt, endsAt, *config, *extrasConfig)
	log.Println("Scheduling playlist...")
	scheduledBlocks := p.ScheduledBlocks()

	if (!skipOutput){
		// output playlist
		log.Println("Outputting playlist...")
		tracks := make([]playlist.XspfTrack, 0)
		for _, scheduleBlock := range scheduledBlocks {
			for _, item := range scheduleBlock.Block.Items {
				tracks = append(tracks, playlist.XspfTrack{Location: "file://" + item.Name()})
			}
		}
		xspf := playlist.XspfPlaylist{Version: "1", Xmlns: "http://xspf.org/ns/0/", XspfTracks: tracks}
		outfile, err := xspf.Output()
		if err != nil {
			log.Fatal("Could not make playlist. %v", err)
		}
		log.Println("Done. Outputted playlist to ", outfile.Name())
	}

	if (!skipPublish){
		// publish playlist
		ok, err := p.Publish(calendarName, scheduledBlocks)
		if( err != nil ){
			log.Fatal("Could not publish playlist.\n%v", err)
		}
		log.Println("Done. Published playlist. ", ok)
	}
	
}