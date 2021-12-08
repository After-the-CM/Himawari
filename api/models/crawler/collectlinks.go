package crawler

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

var tagUrlAttr = map[string][]string{
	"a":       {"href"},
	"applet":  {"code"},
	"area":    {"href"},
	"bgsound": {"src"},
	"body":    {"background"},
	"embed":   {"href", "src"},
	"fig":     {"src"},
	"frame":   {"src"},
	"iframe":  {"src"},
	"img":     {"href", "src", "lowsrc"},
	"input":   {"src"},
	"layer":   {"src"},
	"object":  {"data"},
	"overlay": {"src"},
	"script":  {"src"},
	"table":   {"background"},
	"tb":      {"background"},
	"th":      {"background"},
}

func CollectLinks(body io.Reader, referer *url.URL) {
	doc, err := goquery.NewDocumentFromReader(body)
	if logger.ErrHandle(err) {
		return
	}

	nextStruct := entity.RequestStruct{}
	nextStruct.Referer = referer

	parseHtml(doc, &nextStruct)
	parseForms(doc, &nextStruct)
}

func parseHtml(doc *goquery.Document, r *entity.RequestStruct) {
	for i, v := range tagUrlAttr {
		for _, w := range v {
			doc.Find(i).EachWithBreak(func(_ int, s *goquery.Selection) bool {
				attr, b := s.Attr(w)
				if b {
					var err error
					r.Path, err = url.Parse(attr)
					if logger.ErrHandle(err) {
						return true
					}
					for key, _ := range r.Path.Query() {
						for n, v := range applyData {
							if key == n {
								tmp := r.Path.Query()
								tmp.Del(key)
								tmp.Add(key, v)
								r.Path.RawQuery = tmp.Encode()
							}
						}
					}
					GetRequest(r)
				}
				return true
			})
		}
	}
}

func parseForms(doc *goquery.Document, r *entity.RequestStruct) {
	doc.Find("form").Each(func(_ int, s *goquery.Selection) {
		form := entity.HtmlForm{}
		form.Action, _ = s.Attr("action")
		form.Method, _ = s.Attr("method")
		var inputs []entity.HtmlForm

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			f := form

			typ, ok := s.Attr("type")
			if ok {
				typ = strings.ToLower(typ)
				f.Type = typ
			} else {
				f.Type = "text"
			}

			nameAttr, ok := s.Attr("name")
			if ok {
				f.Name = &nameAttr
			} else {
				f.Name = nil
			}

			value, ok := s.Attr("value")
			if ok {
				f.Value = &value
			} else {
				f.Value = nil
			}

			placeholder, ok := s.Attr("placeholder")
			if ok {
				f.Placeholder = &placeholder
			} else {
				f.Placeholder = nil
			}

			inputs = append(inputs, f)
		})

		s.Find("select").Each(func(_ int, s *goquery.Selection) {
			f := form
			f.IsOption = true
			f.Type = "select"
			nameAttr, ok := s.Attr("name")
			if ok {
				f.Name = &nameAttr
			} else {
				f.Name = nil
			}
			s.Find("option").Each(func(_ int, s *goquery.Selection) {
				value, ok := s.Attr("value")
				if ok {
					f.Options = append(f.Options, value)
				}
			})
			inputs = append(inputs, f)
		})

		s.Find("textarea").Each(func(_ int, s *goquery.Selection) {
			f := form
			f.Type = "textarea"
			f.Value = nil
			nameAttr, ok := s.Attr("name")
			if ok {
				f.Name = &nameAttr
			} else {
				f.Name = nil
			}
			placeholder, ok := s.Attr("placeholder")
			if ok {
				f.Placeholder = &placeholder
			} else {
				f.Placeholder = nil
			}
			if s.Text() != "" {
				text := s.Text()
				f.Placeholder = &text
			}
			inputs = append(inputs, f)
		})
		SetValues(inputs, r)
	})
}
