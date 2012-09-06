package playlist

import (
	"net/http"
	"net/url"
	"errors"
	"io/ioutil"
	"log"
	"bytes"
	"encoding/json"
)

func (p *Playlist) Publish(calendar *Calendar, items []ScheduledBlock) (string, error) {
	log.Printf("Publishing playlist to %s\n", calendar.Name)
	token, err := RefreshAccessToken()
	if (err != nil){
		return "", err
	}
	for _, scheduledBlock := range items {
		if( scheduledBlock.Block.DoPublish ){			
			resp, err := scheduledBlock.Publish(calendar, token)
			if (err != nil){
				return "", err
			}
			var publishResponse interface{}
			err = json.Unmarshal(resp, &publishResponse)
			if (err != nil){
				return "", err
			}
			m := publishResponse.(map[string]interface{})
			if ( m["error"] != nil){
				return "", errors.New(string(resp))
			}
		}
	}
	return "ok", nil
}

func (scheduledBlock ScheduledBlock) Publish(calendar *Calendar, accessToken string) ([]byte, error){
	log.Printf("Publishing %s to %s at %s\n", scheduledBlock.Block.Title, calendar.Name, scheduledBlock.Start.DateTime)
	event := CalendarEvent{ scheduledBlock.Start, scheduledBlock.End, scheduledBlock.Block.Title, "Series: " + scheduledBlock.Block.Series + "\nOriginally aired at: true" }
	return event.Publish(calendar, accessToken)
}

func RefreshAccessToken() (string, error) {
	resp, err := http.PostForm("https://accounts.google.com/o/oauth2/token", url.Values{"refresh_token": {"1/LxxA16-YwWspGM1iXoDqhKFSKNN0BzWm1zkZYKRIQt4"},"client_id": {"404966439763-uv5escoh0lqf7itsp1ifsuvtkjuv9eu9.apps.googleusercontent.com"},"client_secret": {"H7c2qB_YuB6CWKcPg9img88R"},"grant_type": {"refresh_token"}} )
	if (err != nil){
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if (err != nil){
		return "", err
	}
	var accessTokenResponse interface{}
	err = json.Unmarshal(body, &accessTokenResponse)
	if (err != nil){
		return "", err
	}
	m := accessTokenResponse.(map[string]interface{})
	return m["access_token"].(string), nil
}

func (e *CalendarEvent) Publish(calendar *Calendar, accessToken string) ([]byte, error){
	url := "https://www.googleapis.com/calendar/v3/calendars/" + calendar.Id + "/events"
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	br := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", url, br)
	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("Content-Type", "application/json")
	if (err != nil){
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if (err != nil){
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil){
		return nil, err
	}
	return body, nil
}