package main

import (
  "fmt"
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "os"
  "strings"
  "github.com/joho/godotenv"
)

const meetupApiUrl = "https://api.meetup.com/2/members?group_id=%s&key=%s"

type MemberResult struct {
  Results []Member
}

type Member struct {
  Id int
  Name, Bio, Hometown string
  OtherServices struct {
    Twitter struct {
      Identifier string
    }
    Linkedin struct {
      Identifier string
    }
    Facebook struct {
      Identifier string
    }
    Tumblr struct {
      Identifier string
    }
    Flickr struct {
      Identifier string
    }
  } `json:"other_services"`
}

func (m Member) SocialNetworkHandles() []string {
  handles := make([]string, 0)
  if m.OtherServices.Twitter.Identifier != "" {
    twitterHandle := m.OtherServices.Twitter.Identifier
    twitterHandle = strings.Replace(twitterHandle, "@", "", 1)
    handles = append(handles, twitterHandle)
  }
  if m.OtherServices.Linkedin.Identifier != "" {
    linkedinHandle := m.OtherServices.Linkedin.Identifier
    parts := strings.Split(linkedinHandle, "/")
    linkedinHandle = parts[len(parts) - 1]
    handles = append(handles, linkedinHandle)
  }
  if m.OtherServices.Facebook.Identifier != "" {
    facebookHandle := m.OtherServices.Facebook.Identifier
    parts := strings.Split(facebookHandle, "/")
    facebookHandle = parts[len(parts) - 1]
    handles = append(handles, facebookHandle)
  }
  if m.OtherServices.Tumblr.Identifier != "" {
    tumblrHandle := m.OtherServices.Tumblr.Identifier
    // TODO: learn golang replace by regex
    tumblrHandle = strings.Replace(tumblrHandle, "http://", "", 1)
    tumblrHandle = strings.Replace(tumblrHandle, "https://", "", 1)
    parts := strings.Split(tumblrHandle, ".")
    handles = append(handles, parts[0])
  }
  if m.OtherServices.Flickr.Identifier != "" {
    flickrHandle := m.OtherServices.Flickr.Identifier
    parts := strings.Split(flickrHandle, "/")
    lastIndex := len(parts) - 1
    lastPart := parts[lastIndex]
    if lastPart == "" && lastIndex != 0 {
      lastPart = parts[lastIndex - 1]
    }
    handles = append(handles, lastPart)
  }
  return handles
}

func fetchMembers() []Member {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  meetupGroupId := os.Getenv("MEETUP_GROUP_ID")
  meetupApiKey := os.Getenv("MEETUP_API_KEY")
  reqUrl := fmt.Sprintf(meetupApiUrl, meetupGroupId, meetupApiKey)
  client := &http.Client{}
  fmt.Println("API Request URL: ", reqUrl)
  req, reqErr := http.NewRequest("GET", reqUrl, nil)
  if reqErr != nil {
    log.Fatal("NewRequest: ", reqErr)
    return nil
  }
  resp, respErr := client.Do(req)
  if respErr != nil {
    log.Fatal("Do: ", respErr)
    return nil
  }
  defer resp.Body.Close()
  body, dataReadErr := ioutil.ReadAll(resp.Body)
  if dataReadErr != nil {
    log.Fatal("ReadAll: ", dataReadErr)
    return nil
  }
  var results MemberResult
  errr := json.Unmarshal(body, &results)
  if errr != nil {
    log.Fatal(errr)
  }
  return results.Results
}

func main() {
  members := fetchMembers()
  for _, member := range members {
    fmt.Printf("%v\n  %v\n", member.Name, member.Bio)
    for _, handle := range member.SocialNetworkHandles() {
      fmt.Printf("  %v\n", handle)
    }
  }
}
