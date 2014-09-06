package rikimaru

import (
  "net/http"
  "io/ioutil"
  "errors"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/html"
)

type Rikimaru struct {
  url string
  doc *html.HtmlDocument
}

func HttpGet(url string) (htmlString string, err error) {
  resp, err := http.Get(url)

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

func (r *Rikimaru) Init(url string) (err error) {
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
