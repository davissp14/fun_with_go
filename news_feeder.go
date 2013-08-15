package main

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "fun_with_go/rss"
  "fmt"
)

type Channel struct {
  Id          bson.ObjectId `bson:"_id,omitempty"`
  Name        string
  Rss_feed    string
  Description string
}

type NewsItem struct {
  Channel_id  bson.ObjectId
  Title       string
  Link        string
  Comments    string
  Description string
}


func main(){
  session := establishSession("127.0.0.1")
  defer session.Close()

  var channels []Channel

  conn := session.DB("news_feeder_development").C("channels")
  err := conn.Find(bson.M{}).All(&channels)
  if err != nil {
    panic(err)
  }

  for _, channel := range channels {
    fmt.Println("\nPolling for updates on channel: ", channel.Name)
    feed := parseXML(channel.Rss_feed)
    addNewsItems(session, feed, channel)
  }
}

func establishSession(address string) *mgo.Session {
  session, err := mgo.Dial(address)
  if err != nil {
    panic(err)
  }
  return session
}

func parseXML(address string) rss.XMLFeed {
  var rss rss.Rss

  data, err := rss.FetchData(address)
  if err != nil {
    panic(err)
  }

  xmlFeed, err := rss.ToXml(data)
  if err != nil {
    panic(err)
  }
  return xmlFeed
}

func addNewsItems(session *mgo.Session, feed rss.XMLFeed, channel Channel) {
  collection := session.DB("news_feeder_development").C("news_items")
  added := 0
  for _, item := range feed.XMLFeed.XMLChannel {
    var newsItem NewsItem
    newsItem = NewsItem{Title: item.Title, Link: item.Link, Comments: item.Comments, Description: item.Description, Channel_id: channel.Id}
    err := collection.Insert(&newsItem)
    if err == nil {
      added++
      fmt.Println("Adding title: " + item.Title)
    }
  }
  fmt.Println(added, " titles added!")
}


