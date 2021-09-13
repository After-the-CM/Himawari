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
	parseHtml(doc, nextStruct)
	// ParseForms中でfunc3を呼ぶ
	parseForms(doc, nextStruct)
}

// paramerを必要としないタグからリンクを収集 → func1に投げる
// ※ この実装では絶対うまく行かない。nextStructの初期化タイミング、Returnされるタイミングを丁寧に考える必要がある。
// ValidateOrigin(checkorigin)とIsExistをfunc2に戻さないと実装はきつそう(?)。自分がアルゴリズム弱太郎なだけであってほしい。
func parseHtml(doc *goquery.Document, r entity.RequestStruct) {

	tagUrlAttr := map[string][]string{
		"a":       {"href"},
		"applet":  {"code"},
		"area":    {"href"},
		"bgsound": {"src"},
		"body":    {"background"},
		"embed":   []string{"href", "src"},
		"fig":     {"src"},
		"frame":   {"src"},
		"iframe":  {"src"},
		"img":     []string{"href", "src", "lowsrc"},
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

//github.com/jay/wget/blob/099d8ee3da3a6eea5635581ae517035165f400a5/src/html-url.c
// url.Parseはgetrequest.go, postrequest.goで実装されている

func parseForms(doc *goquery.Document, r entity.RequestStruct) (forms []entity.HtmlForm) {
	doc.Find("form").Each(func(_ int, s *goquery.Selection) {
		form := entity.HtmlForm{Values: url.Values{}}
		form.Action, _ = s.Attr("action")
		form.Method, _ = s.Attr("method")

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			tag := "input"

			name, nameB := s.Attr("name")
			if nameB {
				form.Values.Add("Name", name)
			} else {
				// fmt.Fprintln(os.Stderr, "'name' Not Found...")
				form.Values.Add("Name", "NaN")
			}

			typ, typB := s.Attr("type")
			typ = strings.ToLower(typ)
			if typB {
				form.Values.Add("Type", typ)
			} else {
				// fmt.Fprintln(os.Stderr, "'type' Not Found...")
				form.Values.Add("Type", "NaN")
			}

			_, checked := s.Attr("checked")
			if (typ == "radio" || typ == "checkbox") && !checked {
				//return
			}

			form.Values.Add("Tag", tag)

			value, valueB := s.Attr("value")
			if valueB {
				form.Values.Add("Value", value)
			} else {
				// fmt.Fprintln(os.Stderr, "'value' Not Found...")
				form.Values.Add("Value", "NaN")
			}

			placeholder, placeholderB := s.Attr("placeholder")
			if placeholderB {
				form.Values.Add("Placeholder", placeholder)
			} else {
				form.Values.Add("Placeholder", "NaN")
			}

			pattern, patternB := s.Attr("pattern")
			if patternB {
				form.Values.Add("Pattern", pattern)
			} else {
				// fmt.Fprintln(os.Stderr, "'value' Not Found...")
				form.Values.Add("Pattern", "NaN")
			}

			require, requireB := s.Attr("require")
			if requireB {
				form.Values.Add("Require", require)
			} else {
				form.Values.Add("Require", "Nan")
			}
		})

		SetValues(form, r)

		/*
			s.Find("textarea").Each(func(_ int, s *goquery.Selection) {
				name, _ := s.Attr("name")
				if name == "" {
					//return
				}

				value := s.Text()
				form.Values.Add(name, value)
			})
		*/

		// formタグがある限り収集したParameterとその値はformsにて保持できている。
		// fmt.Println(fomrs)
		forms = append(forms, form)

	})

	return forms
}
