package sitemap

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"Himawari/models/entity"
)

func addChild(node *entity.Node, parsedPath []string, request http.Request, time float64) {
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
		addChild(&(*node.Children)[childIdx], parsedPath[1:], request, time)
	} else {
		fmt.Println(request)
		if !isExist(node, []string{}, request) {
			(*node).Messages = append((*node).Messages, entity.Message{
				Request: request,
				Time:    time,
			})
		}
	}

}

func Add(request http.Request, time float64) {
	parsedPath := strings.Split(request.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)

	// paramに `;` があるとクエリのパースでバグるため、`;` だけURLエンコード
	request.URL.RawQuery = strings.Replace(request.URL.RawQuery, ";", "%3B", -1)
	request.PostForm, _ = url.ParseQuery(strings.Replace(request.PostForm.Encode(), ";", "%3B", -1))

	addChild(&entity.Nodes, parsedPath, request, time)
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

func IsExist(request http.Request) bool {
	parsedPath := strings.Split(request.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)

	request.URL.RawQuery = strings.Replace(request.URL.RawQuery, ";", "%3B", -1)
	request.PostForm, _ = url.ParseQuery(strings.Replace(request.PostForm.Encode(), ";", "%3B", -1))

	return isExist(&entity.Nodes, parsedPath, request)
}

func isExist(node *entity.Node, parsedPath []string, request http.Request) bool {
	if len(parsedPath) > 0 {
		childIdx := getChildIdx(node, parsedPath[0])

		if childIdx >= 0 {
			return isExist(&(*node.Children)[childIdx], parsedPath[1:], request)
		} else {
			return false
		}
	} else {

		for _, msg := range node.Messages {
			if msg.Request.URL.RawQuery == request.URL.RawQuery && msg.Request.PostForm.Encode() == request.PostForm.Encode() {
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
