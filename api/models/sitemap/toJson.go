package sitemap

import (
	"net/http/cookiejar"
	"net/url"
	"path"
	"sort"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

func messages(node *entity.Node) []entity.JsonMessage {
	msg := make([]entity.JsonMessage, len(node.Messages))
	for i, m := range node.Messages {
		msg[i] = entity.JsonMessage{
			Time:       m.Time,
			Referer:    m.Request.Referer(),
			GetParams:  m.Request.URL.Query(),
			PostParams: m.Request.PostForm,
			URL:        m.Request.URL.Scheme + "://" + m.Request.URL.Host + m.Request.URL.Path,
		}
	}
	return msg
}

func jsonAddChild(node *entity.Node, jsonNode *entity.JsonNode, url *url.URL, jar *cookiejar.Jar) {
	if node.Children != nil {
		for i, n := range *node.Children {
			// url.URLをstringにしてurl.URLにパースすることで元のurlを書き換えないように
			u, err := url.Parse(url.String())
			logger.ErrHandle(err)
			u.Path = path.Join(u.Path, n.Path)
			child := &entity.JsonNode{
				Path:     n.Path,
				Cookies:  getCookies(u, jar),
				Messages: messages(&n),
				URL:      u.String(),
			}

			(*jsonNode).Children = append((*jsonNode).Children, *child)
			if node.Children != nil && *node.Children != nil {
				childURL, err := url.Parse(child.URL)
				logger.ErrHandle(err)
				jsonAddChild(&(*node.Children)[i], &jsonNode.Children[i], childURL, jar)
			}
		}
	}
}

func Merge(url *url.URL, jar *cookiejar.Jar) {
	entity.JsonNodes = entity.JsonNode{
		Path:     entity.Nodes.Path,
		Cookies:  getCookies(url, jar),
		Messages: messages(&entity.Nodes),
		URL:      url.Scheme + "://" + url.Host,
	}
	jsonAddChild(&entity.Nodes, &entity.JsonNodes, url, jar)
}

func getCookies(url *url.URL, jar *cookiejar.Jar) []entity.JsonCookie {
	jarcookies := jar.Cookies(url)
	cookies := make([]entity.JsonCookie, len(jarcookies))
	for i, jarcookie := range jarcookies {
		cookies[i].Path = jarcookie.Path
		cookies[i].Name = jarcookie.Name
		cookies[i].Value = jarcookie.Value
	}
	return cookies
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

func InitVulnMap(node *entity.JsonNode) {
	for i := 0; len(node.Children) > i; i++ {
		InitVulnMap(&node.Children[i])
	}
	for _, v := range node.Issue {
		entity.Vulnmap[v.Kind].Issues = append(entity.Vulnmap[v.Kind].Issues, v)
	}
}
