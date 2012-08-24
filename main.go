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

func main() {
	flag.StringVar(&franchiseName, "franchise", "StarCraft 2", "Name of franchise. Default is StarCraf 2.")
	flag.StringVar(&startsAtTime, "startsAt", time.Now().Format(timeFormat), "Start time. Default is now.")
	flag.Parse() // parses the flags

	startsAt, err := time.Parse(timeFormat, startsAtTime)
    if err != nil {
        fmt.Println(err)
        return
    }	

	fmt.Printf("\n* Starting...\n")

	fmt.Printf("Franchise\n\targ: %s\n", franchiseName)

	fmt.Printf("Starts\n\targ: %v\n\ttime: %v\n", startsAtTime, startsAt)

	fmt.Printf("Done.\n")
}	