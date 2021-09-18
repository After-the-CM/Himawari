package crawler

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"Himawari/models/entity"

	"github.com/PuerkitoBio/goquery"
)

//func2
func CollectLinks(body io.Reader, referer *url.URL) {
	fmt.Println("Start func2")
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	title := doc.Find("title").Text()
	fmt.Println("ページタイトル：" + title)

	nextStruct := entity.RequestStruct{}
	nextStruct.Referer = referer

	// (ちゃんと戻ってこれる？心配。)
	parseHtml(doc, &nextStruct)
	// ParseForms中でfunc3を呼ぶ
	parseForms(doc, &nextStruct)
}

// paramerを必要としないタグからリンクを収集 → func1に投げる
// ※ この実装では絶対うまく行かない。nextStructの初期化タイミング、Returnされるタイミングを丁寧に考える必要がある。
// ValidateOrigin(checkorigin)とIsExistをfunc2に戻さないと実装はきつそう(?)。自分がアルゴリズム弱太郎なだけであってほしい。
func parseHtml(doc *goquery.Document, r *entity.RequestStruct) {

	tagUrlAttr := map[string][]string{
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

	for i, v := range tagUrlAttr {
		for _, w := range v {
			doc.Find(i).Each(func(_ int, s *goquery.Selection) {
				attr, b := s.Attr(w)
				if b {
					r.Path, _ = url.Parse(attr)
					GetRequest(r)
				}
			})
		}
	}

}

func parseForms(doc *goquery.Document, r *entity.RequestStruct) {
	doc.Find("form").Each(func(_ int, s *goquery.Selection) {
		form := entity.HtmlForm{}
		form.Action, _ = s.Attr("action")

		form.Method, _ = s.Attr("method")
		form.Method = strings.ToUpper(form.Method)
		var inputs []entity.HtmlForm

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			f := form
			//onclickがあるinputはreturnしている。
			/*
				_, ok := s.Attr("onclick")
				if ok {
					return
				}
			*/

			typ, ok := s.Attr("type")
			if ok {
				typ = strings.ToLower(typ)
				f.Type = typ
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
			//textareaタグにはvalue属性がないため
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
				//placeholderの優先度はTestDataよりも高いためplaceholderに入れておく
				text := s.Text()
				f.Placeholder = &text
			}
			inputs = append(inputs, f)
		})

		SetValues(inputs, r)

	})

}
