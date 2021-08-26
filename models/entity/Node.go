package entity

type Node struct {
	Parent   *Node
	Path     string
	Children *[]Node
}

type JsonNode struct {
	Path     string     `json:"path"`
	Children []JsonNode `json:"children"`
}

var Nodes = Node{
	Path: "/",
}
