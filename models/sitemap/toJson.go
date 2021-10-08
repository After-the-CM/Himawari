package sitemap

import (
	"sort"

	"Himawari/models/entity"
)

func messages(node *entity.Node) []entity.JsonMessage {
	msg := make([]entity.JsonMessage, len(node.Messages))
	for i, m := range node.Messages {
		msg[i] = entity.JsonMessage{
			Time:       m.Time,
			Referer:    m.Request.Referer(),
			GetParams:  m.Request.URL.Query(),
			PostParams: m.Request.PostForm,
		}
	}
	return msg
}

func jsonAddChild(node *entity.Node, jsonNode *entity.JsonNode) {
	if node.Children != nil {
		for i, n := range *node.Children {
			url := n.Messages[0].Request.URL
			child := &entity.JsonNode{
				Path:     n.Path,
				URL:      url.Scheme + "://" + url.Host + url.Path,
				Messages: messages(&n),
			}

			(*jsonNode).Children = append((*jsonNode).Children, *child)
			if node.Children != nil && *node.Children != nil {
				jsonAddChild(&(*node.Children)[i], &jsonNode.Children[i])
			}
		}
	}
}

func Merge(url string) {
	entity.JsonNodes = entity.JsonNode{
		Path:     entity.Nodes.Path,
		URL:      url,
		Messages: messages(&entity.Nodes),
	}
	jsonAddChild(&entity.Nodes, &entity.JsonNodes)
}

func SortJson() {
	sortChild(entity.JsonNodes)
}

func sortChild(node entity.JsonNode) {
	if node.Children != nil {
		sort.Sort(node)
		for _, child := range node.Children {
			sortChild(child)
		}
	}
}
