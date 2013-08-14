package main

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/xml"
  "time"
)

type Channel struct {
  Id          bson.ObjectId `bson:"_id,omitempty"`
  Name        string
  Rss_feed    string
  Description string
}

type NewsItem struct {
  Id          bson.ObjectId `bson:"_id,omitempty"`
  Name        string
  Title       string
  Link        string
  Comments    string
  Description string
}

type XMLItems struct {
  XMLName     xml.Name `xml:"item"`
  Title       string   `xml:"title"`
  Channel_id  bson.ObjectId
  Link        string   `xml:"link"`
  Comments    string   `xml:"comments"`
  Description string   `xml:"description"`
}

type XMLChannel struct {
  XMLName     xml.Name `xml:"channel"`
  Title       string   `xml:"title"`
  Link        string   `xml:"link"`
  Description string   `xml:"description"`
  XMLChannel []XMLItems   `xml:"item"`
}

type XMLFeed struct {
  XMLName xml.Name `xml:"rss"`
  XMLFeed XMLChannel `xml:"channel"`
}

func main(){
  session := establishSession("127.0.0.1")
  defer session.Close()

  result := Channel{}

  c := session.DB("news_feeder_development").C("channels")
  err := c.Find(bson.M{"name": "Hacker News"}).One(&result)
  if err != nil {
    panic(err)
  }
  channel_id := result.Id

  ticker := time.NewTicker(time.Second * 300)
  go func() {
    for t := range ticker.C {
      fmt.Println("Running update at: ", t)
      var parsedFeed XMLFeed

      data := retrieveXmlData(result.Rss_feed)
      xml.Unmarshal(data, &parsedFeed)

      insertNewsItems(session, parsedFeed, channel_id)
    }
  }()

  time.Sleep(time.Second * 3600)
  ticker.Stop()
  fmt.Println("ticket has been stopped")
}

func establishSession(address string) *mgo.Session {
  session, err := mgo.Dial(address)
  if err != nil {
    panic(err)
  }
  return session
}

func retrieveXmlData(address string) []byte {
  resp, err := http.Get(address)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  data, _ := ioutil.ReadAll(resp.Body)

  return data
}

func insertNewsItems(session *mgo.Session, parsedFeed XMLFeed, channel_id bson.ObjectId) {
  i := session.DB("news_feeder_development").C("news_items")
  for _, item := range parsedFeed.XMLFeed.XMLChannel {
    item.Channel_id = channel_id
    err := i.Insert(&item)
    if err != nil {
      fmt.Println("Duplicate found, ignoring title: " + item.Title)
    } else {
      fmt.Println("Adding title: " + item.Title)
    }
  }
}

