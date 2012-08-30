package playlist

import (
	"os"
	"log"
	"encoding/xml"
)

type XspfPlaylist struct {
    XMLName xml.Name `xml:"playlist"`
    Version string `xml:"version,attr"`
    Xmlns string `xml:"xmlns,attr"`
    XspfTracks []XspfTrack `xml:"trackList>track"`
}

type XspfTrack struct {
	Location string `xml:"location"`
}

func (xspfPlaylist *XspfPlaylist) Output() (*os.File, error) {
	xmlstring, err := xml.MarshalIndent(xspfPlaylist, "", "    ")
	if err != nil {
	    log.Fatalf("xml.MarshalIndent: %v", err)
	}
	output, err := os.Create("out.xspf")
	if err != nil {
	    return nil, err
	}
	_, err = output.Write( []byte(xml.Header + string(xmlstring)) )
	if err != nil {
	    return nil, err
	}
	err = output.Close()
	if err != nil {
	    return nil, err
	}
	return output, nil
}