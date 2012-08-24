package main

import (
	"fmt"
	"flag"
	"time"
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

func init() {
	now := time.Now()
	flag.StringVar(&franchiseName, "franchise", "StarCraft 2", "Name of franchise. Default is StarCraf 2.")
	flag.StringVar(&startsAtTime, "start", now.Format(timeFormat), "Start time. Default is now.")
	flag.StringVar(&endsAtTime, "end", now.Add(time.Hour * 24).Format(timeFormat), "End time. Default is 24 hours from now.")
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

	// flags
	fmt.Printf("Franchise\n\targ: %s\n", franchiseName)
	fmt.Printf("Starts\n\targ: %v\n\ttime: %v\n", startsAtTime, startsAt)
	fmt.Printf("Ends\n\targ: %v\n\ttime: %v\n", endsAtTime, endsAt)

	fmt.Println("Done.")
}