package sitemap

import (
	"fmt"
	"net/http"
	"strings"

	"Himawari/models/entity"
)

func addChild(node *entity.Node, parsedPath []string, request http.Request) {
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
		if len(parsedPath) > 1 {
			addChild(&(*node.Children)[childIdx], parsedPath[1:], request)
		} else {
			for _, v := range (*node.Children)[childIdx].Messages {
				if v.URL.RawQuery == request.URL.RawQuery {
					return
				}
			}
			(*node.Children)[childIdx].Messages = append((*node.Children)[childIdx].Messages, request)
		}
	}
}

func AddPath(node *entity.Node, request http.Request) {
	parsedPath := strings.Split(request.URL.Path, "/")
	parsedPath = removeSpace(parsedPath)
	addChild(node, parsedPath, request)
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
		for _, v := range node.Messages {
			if v.URL.Query().Encode() == request.URL.Query().Encode() {
				return true
			}
		}
		return false
	}
}

func jsonAddChild(node entity.Node, jsonNode *entity.JsonNode) {
	if node.Children != nil {
		for i, v := range *node.Children {
			child := &entity.JsonNode{
				Path: v.Path,
			}

			(*jsonNode).Children = append((*jsonNode).Children, *child)
			if v.Children != nil && *v.Children != nil {
				jsonAddChild(v, &jsonNode.Children[i])
			}
		}
	}
}

func MtoJ(node entity.Node) entity.JsonNode {
	jsonNode := entity.JsonNode{
		Path: node.Path,
	}
	jsonAddChild(node, &jsonNode)
	return jsonNode
}

func PrintMap(node entity.Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("\t")
	}
	fmt.Println(node.Path)

	for i := 1; i < len(node.Messages); i++ {
		fmt.Printf("%v, ", node.Messages[i].URL.Query().Encode())
	}

	fmt.Println()
	if node.Children != nil {
		indent++
		for _, v := range *node.Children {
			PrintMap(v, indent)
		}
	}
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
