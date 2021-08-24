package sitemap

import (
	"fmt"
	"strings"

	"mbsd/models/entity"
)

func addChild(node *entity.Node, parsedPath []string) {
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
			addChild(&(*node.Children)[childIdx], parsedPath[1:])
		}
	}
}

func AddPath(node *entity.Node, fullPath string) {
	parsedPath := strings.Split(fullPath, "/")
	parsedPath = removeSpace(parsedPath)
	addChild(node, parsedPath)
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
