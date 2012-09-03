package playlist

import (
	"fmt"
	//"net/http"
	//"net/url"
	//"strings"
	//"io/ioutil"
	"log"
	//"bytes"
	//"time"
	//"encoding/json"
)

func (p *Playlist) Publish(calendarName string, items []ScheduledBlock) () {
	log.Printf("Publishing playlist to %s\n", calendarName)
	token, err := RefreshAccessToken()
	if (err != nil){
		log.Fatal("Unable to obtain access token ", err)
	}
	fmt.Printf("token: %s\n", token)
	for _, scheduledBlock := range items {
		if( scheduledBlock.Block.DoPublish ){			
			scheduledBlock.Publish(calendarName, token)
		}
	}
	return
}

func (scheduledBlock ScheduledBlock) Publish(calendarName string, accessToken string) (string){
	log.Printf("Publishing %s to %s at %s\n", scheduledBlock.Block.Title, calendarName, scheduledBlock.Start.DateTime)
	//event := CalendarEvent{}
	return accessToken
}

func RefreshAccessToken() (string, error) {
	// resp, err := http.PostForm("https://accounts.google.com/o/oauth2/token", url.Values{"refresh_token": {"1/LxxA16-YwWspGM1iXoDqhKFSKNN0BzWm1zkZYKRIQt4"},"client_id": {"404966439763-uv5escoh0lqf7itsp1ifsuvtkjuv9eu9.apps.googleusercontent.com"},"client_secret": {"H7c2qB_YuB6CWKcPg9img88R"},"grant_type": {"refresh_token"}} )
	// if (err != nil){
	// 	return "", err
	// }
	// body, err := ioutil.ReadAll(resp.Body)
	// defer resp.Body.Close()
	// if (err != nil){
	// 	return "", err
	// }
	// var accessTokenResponse interface{}
	// err = json.Unmarshal(body, &accessTokenResponse)
	// if (err != nil){
	// 	return "", err
	// }
	// m := accessTokenResponse.(map[string]interface{})
	// return m["access_token"].(string), nil

	return "ya29.AHES6ZQvJBnW3AJ75a7TBuCQlUfZwcH0R5XIwbP-86Mx84GRkk1jAzI", nil
}

// func (e *CalendarEvent) Publish(accessToken string) (string, error){
// 	url := "https://www.googleapis.com/calendar/v3/calendars/" + e.Calendar.Id + "/events"

	//execute( parse_url("https://www.googleapis.com/calendar/v3/calendars/:calendarId/events", parameters), data, { :verb => "post" } )
	// b, err := json.Marshal(`{}`)
	// if err != nil {
	// 	log.Fatal("JSON Error: ", err)
	// }
	// br := bytes.NewBuffer(b)
	// log.Println("b: ", br)

	// req, err := http.NewRequest("POST", "https://accounts.google.com/o/oauth2/token", br)
	// req.Header.Add("Content-Type", `"application/json"`)
	// if (err != nil){
	// 	log.Fatalf("http.NewRequest POST error: %v", err)
	// }
	// resp, err := http.DefaultClient.Do(req)
	// if (err != nil){
	// 	log.Fatalf("client.Do error: %v", err)
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if (err != nil){
	// 	log.Fatalf("ioutil.ReadAll error: %v", err)
	// }
	// fmt.Printf("resp body: %s\n", body)

	//return url, nil
//}