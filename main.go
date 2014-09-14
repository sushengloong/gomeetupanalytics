package main

import (
  "fmt"
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "os"
  "github.com/joho/godotenv"
)

const meetupApiUrl = "https://api.meetup.com/2/members?group_id=%s&key=%s"

type MemberResult struct {
  Results []Member
}

type Member struct {
  Id int
  Name, Bio string
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
    fmt.Printf("%v: %v (%v)\n", member.Id, member.Name, member.Bio)
  }
}
