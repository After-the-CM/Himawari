package sitemap

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"Himawari/models/entity"
)

func addChild(node *entity.Node, parsedPath []string, req http.Request, time float64) {
	if len(parsedPath) > 0 {
		childIdx := getChildIdx(node, parsedPath[0])
		child := &entity.Node{
			Parent: node,
			Path:   parsedPath[0],
		}
		if childIdx == -1 {
			// pathがchildrenにない
			*node.Children = append(*node.Children, *child)
			childIdx = len(*node.Children) - 1
		} else if childIdx == -2 {
			children := make([]entity.Node, 0)
			node.Children = &children
			*node.Children = append(*node.Children, *child)
			childIdx = 0
		}
		addChild(&(*node.Children)[childIdx], parsedPath[1:], req, time)
	} else {
		if !isExist(node, []string{}, req) {
			(*node).Messages = append((*node).Messages, entity.Message{
				Request: req,
				Time:    time,
			})
		}
	}

}

func Add(req http.Request, time float64) {
	parsedPath := strings.Split(req.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)

	// paramに `;` があるとクエリのパースでバグるため、`;` だけURLエンコード
	req.URL.RawQuery = strings.Replace(req.URL.RawQuery, ";", "%3B", -1)
	req.PostForm, _ = url.ParseQuery(strings.Replace(req.PostForm.Encode(), ";", "%3B", -1))

	addChild(&entity.Nodes, parsedPath, req, time)
}

func getChildIdx(node *entity.Node, path string) int {
	if (*node).Children != nil {
		for i, child := range *node.Children {
			if child.Path == path {
				return i
			}
		}
		return -1
	}
	return -2
}

func IsExist(req http.Request) bool {
	parsedPath := strings.Split(req.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)

	req.URL.RawQuery = strings.Replace(req.URL.RawQuery, ";", "%3B", -1)
	req.PostForm, _ = url.ParseQuery(strings.Replace(req.PostForm.Encode(), ";", "%3B", -1))

	return isExist(&entity.Nodes, parsedPath, req)
}

func isExist(node *entity.Node, parsedPath []string, req http.Request) bool {
	if len(parsedPath) > 0 {
		childIdx := getChildIdx(node, parsedPath[0])

		if childIdx >= 0 {
			return isExist(&(*node.Children)[childIdx], parsedPath[1:], req)
		} else {
			return false
		}
	} else {

		for _, msg := range node.Messages {
			if msg.Request.URL.RawQuery == req.URL.RawQuery && msg.Request.PostForm.Encode() == req.PostForm.Encode() {
				return true
			}
		}
		return false
	}
}

func printMap(node entity.Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("\t")
	}
	fmt.Println(node.Path)
	for i := 0; i < len(node.Messages); i++ {
		fmt.Printf("%v, ", node.Messages[i])
		fmt.Println()
	}
	if node.Children != nil {
		indent++
		for _, v := range *node.Children {
			printMap(v, indent)
		}
	}
}

func PrintMap() {
	printMap(entity.Nodes, 0)
}

func removeSpace(parsedPath []string) []string {
	paths := make([]string, 0)
	for _, v := range parsedPath {
		if v == "" {
			continue
		}
		paths = append(paths, v)
	}
	return paths
}

func Reset() {
	entity.Nodes = entity.Node{
		Path: "/",
	}

	entity.JsonNodes = entity.JsonNode{}
}
