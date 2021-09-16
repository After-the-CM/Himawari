package entity

import (
	"net/http"
	"net/url"
)

type Node struct {
	Parent   *Node
	Path     string
	Children *[]Node
	Messages []Message
}

type Message struct {
	Request http.Request
	Time    float64
}

type JsonNode struct {
	Path     string        `json:"path"`
	Messages []JsonMessage `json:"messages"`
	Children []JsonNode    `json:"children"`
}

type JsonMessage struct {
	Times      float64    `json:"times"`
	GetParams  url.Values `json:"getParams"`
	PostParams url.Values `json:"postParams"`
}

var Nodes = Node{
	Path: "/",
}

var JsonNodes = JsonNode{}
