package playlist

import (
	"time"
	"os"
	"encoding/xml"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

// Playlist is a description of a set of media files.
type Playlist struct {
	StartsAt, EndsAt time.Time
	Config yaml.File
	Items []*AvailableBlock
	ExtrasConfig yaml.File
	ExtraItems []*AvailableBlock
}

// An AvailableBlock is a set of related media files available for scheduling.
type AvailableBlock struct {
	Title string
	Series string
	Items []*os.File
	Duration int
	DoPublish bool
}

type ScheduledBlock struct {
	Block *AvailableBlock
	Start CalendarTime
 	End CalendarTime
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

type Calendar struct {
	Id string
	Name string
}

type CalendarTime struct {
	DateTime time.Time `json:"dateTime"`
}

type CalendarEvent struct {
	Start CalendarTime `json:"start"`
	End CalendarTime `json:"end"`
	Summary string `json:"summary"`
	Description string `json:"description"`
}