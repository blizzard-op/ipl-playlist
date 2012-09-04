package playlist

import (
	"os"
	"log"
	"encoding/xml"
)

func (xspfPlaylist *XspfPlaylist) Output(filepath string) (*os.File, error) {
	xmlstring, err := xml.MarshalIndent(xspfPlaylist, "", "    ")
	if err != nil {
	    log.Fatalf("xml.MarshalIndent: %v", err)
	}
	output, err := os.Create(filepath)
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