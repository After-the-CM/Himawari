package sitemap

import (
	"Himawari/models/entity"
)

func messages(node *entity.Node) []entity.JsonMessage {
	msg := make([]entity.JsonMessage, len(node.Messages))
	for i, m := range node.Messages {
		msg[i].Times = m.Time
		msg[i].GetParams = m.Request.URL.Query().Encode()
		msg[i].PostParams = m.Request.PostForm.Encode()
	}
	return msg
}

func jsonAddChild(node *entity.Node, jsonNode *entity.JsonNode) {
	if node.Children != nil {
		for i, n := range *node.Children {
			child := &entity.JsonNode{
				Path:     n.Path,
				Messages: messages(&n),
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
		Path:     entity.Nodes.Path,
		Messages: messages(&entity.Nodes),
	}
	jsonAddChild(&entity.Nodes, &jsonNode)
	return jsonNode
}
