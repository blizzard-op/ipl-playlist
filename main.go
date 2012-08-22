package main

import (
	"fmt"
	"flag"
)

var franchiseName string;

func init() {
	fmt.Printf("Initializing...\n")
	flag.StringVar(&franchiseName, "franchise", "StarCraft 2", "name of franchise")
}

func main() {
	flag.Parse() // parses the logging flags

	fmt.Printf("Starting...\n")

	fmt.Printf("n %s\n", franchiseName)

	fmt.Printf("Done.\n")
}	