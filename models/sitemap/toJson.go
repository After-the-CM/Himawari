package sitemap

import (
	"Himawari/models/entity"
	"fmt"
	"os"
)

func getParams(node *entity.Node) []string {
	params := make([]string, len((*node).Messages))
	for i, message := range (*node).Messages {
		params[i] = message.URL.Query().Encode()
	}
	return params
}

func postParams(node *entity.Node) []string {
	params := make([]string, len((*node).Messages))
	for i, message := range (*node).Messages {
		err := message.ParseForm()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		params[i] = message.PostForm.Encode()
	}
	return params
}

func jsonAddChild(node *entity.Node, jsonNode *entity.JsonNode) {
	if node.Children != nil {
		for i, v := range *node.Children {
			child := &entity.JsonNode{
				Path:       v.Path,
				GetParams:  getParams(&(*node.Children)[i]),
				PostParams: postParams(&(*node.Children)[i]),
			}

			(*jsonNode).Children = append((*jsonNode).Children, *child)
			if node.Children != nil && *node.Children != nil {
				jsonAddChild(&(*node.Children)[i], &jsonNode.Children[i])
			}
		}
	}
}

func Json() entity.JsonNode {
	jsonNode := entity.JsonNode{
		Path:       entity.Nodes.Path,
		GetParams:  getParams(&entity.Nodes),
		PostParams: postParams(&entity.Nodes),
	}
	jsonAddChild(&entity.Nodes, &jsonNode)
	return jsonNode
}
