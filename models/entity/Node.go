package entity

import "net/http"

type Node struct {
	Parent   *Node
	Path     string
	Children *[]Node
	Messages []http.Request
}

type JsonNode struct {
	Path     string     `json:"path"`
	Children []JsonNode `json:"children"`
}

var Nodes = Node{
	Path: "/",
}
