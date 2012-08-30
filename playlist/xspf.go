package playlist

import "encoding/xml"

type XspfPlaylist struct {
    XMLName xml.Name `xml:"playlist"`
    Version string `xml:"version,attr"`
    Xmlns string `xml:"xmlns,attr"`
    XspfTracks []XspfTrack `xml:"trackList>track"`
}

type XspfTrack struct {
	Location string `xml:"location"`
}