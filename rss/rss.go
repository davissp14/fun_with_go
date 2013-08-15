package rss

import (
  "net/http"
  "io/ioutil"
  "encoding/xml"
)

type XMLItems struct {
  XMLName     xml.Name `xml:"item"`
  Title       string   `xml:"title"`
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

type Rss []interface{}

func (rss Rss) FetchData(address string) ([]byte, error) {
  resp, err := http.Get(address)
  defer resp.Body.Close()
  data, _ := ioutil.ReadAll(resp.Body)

  return data, err
}

func (rss Rss) ToXml(data []byte) (XMLFeed, error) {
  var feed XMLFeed
  err := xml.Unmarshal(data, &feed)

  return feed, err
}


