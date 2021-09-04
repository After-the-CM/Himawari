package sitemap

import (
	"fmt"
	"net/http"
	"strings"

	"Himawari/models/entity"
)

func addChild(node *entity.Node, parsedPath []string, request *http.Request) {
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
			// children == nil
			children := make([]entity.Node, 0)
			node.Children = &children
			*node.Children = append(*node.Children, *child)
			childIdx = 0
		}
		addChild(&(*node.Children)[childIdx], parsedPath[1:], request)
	} else {
		if !isExist(node, []string{}, request) {
			(*node).Messages = append((*node).Messages, request)
		}
	}

}

func Add(request *http.Request) {
	parsedPath := strings.Split(request.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)
	addChild(&entity.Nodes, parsedPath, request)
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

func IsExist(request *http.Request) bool {
	parsedPath := strings.Split(request.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)
	return isExist(&entity.Nodes, parsedPath, request)
}

func isExist(node *entity.Node, parsedPath []string, request *http.Request) bool {
	if len(parsedPath) > 0 {
		childIdx := getChildIdx(node, parsedPath[0])

		if childIdx >= 0 {
			return isExist(&(*node.Children)[childIdx], parsedPath[1:], request)
		} else {
			return false
		}
	} else {
		for _, msg := range node.Messages {
			if msg.URL.RawQuery == request.URL.RawQuery && msg.PostForm.Encode() == request.PostForm.Encode() {
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
