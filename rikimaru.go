package rikimaru

import (
  "net/http"
  "io/ioutil"
  "errors"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/html"
  "time"
)

type Rikimaru struct {
  url string
  doc *html.HtmlDocument
}

func HttpGet(url string) (htmlString string, err error) {
  // resp, err := http.Get(url)
  timeout := time.Duration(5 * time.Second)
  client := &http.Client{
    Timeout: timeout,
  }

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return "", err
  }

  req.Header.Set("Referer", "https://google.com/")
  req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.124 Safari/537.36")

  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }

  html, err := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  if err != nil {
    return "", err
  }

  return string(html), err
}

func (r *Rikimaru) Doc() *html.HtmlDocument {
  return r.doc
}

func (r *Rikimaru) InitWithURL(url string) (err error) {
  r.url = url
  if len(r.url) == 0 {
    err := errors.New("URL not set")

    return err
  }

  html, err := HttpGet(r.url)

  if err != nil {
    return err
  }

  r.doc, err = gokogiri.ParseHtml([]byte(html))

  return err
}

func (r *Rikimaru) InitWithText(text string) (err error) {
  r.doc, err = gokogiri.ParseHtml([]byte(text))

  return err
}

func (r *Rikimaru) AddNodeAttribute(xpathString, attributeName, attributeValue string) (err error) {
  nodes, err := r.doc.Search(xpathString)

  if err != nil {
    return err
  }

  if len(nodes) < 1 {
    return errors.New("No node found")
  }

  node := nodes[0]

  if node.Attribute(attributeName) == nil {
    node.SetAttr(attributeName, attributeValue)
  } else {
    node.SetAttr(attributeName, node.Attribute(attributeName).String() + " " + attributeValue)
  }

  return err
}

func (r *Rikimaru) Remove(xpathString string) (err error) {
  nodes, err := r.doc.Search(xpathString)

  if err != nil {
    return err
  }

  for _, node := range nodes {
    node.Remove()
  }

  return err
}

func (r *Rikimaru) SaveToFile(filename string) (err error) {
  err = ioutil.WriteFile(filename, []byte(r.doc.String()), 0644)

  return err
}
