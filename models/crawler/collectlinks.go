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

	parseHtml(doc, &nextStruct)
	parseForms(doc, &nextStruct)
}

func parseHtml(doc *goquery.Document, r *entity.RequestStruct) {
	// https://github.com/jay/wget/blob/099d8ee3da3a6eea5635581ae517035165f400a5/src/html-url.c
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		//pathを*url.URL型に変更
		//r.Path = &href
		r.Path, _ = url.Parse(href)

		// url.Parseはgetrequest.go, postrequest.goで実装されている
		GetRequest(r)
	})
	doc.Find("applet").Each(func(_ int, s *goquery.Selection) {
		code, _ := s.Attr("code")
		//r.Path = &code
		r.Path, _ = url.Parse(code)

		GetRequest(r)
	})
	doc.Find("area").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		//r.Path = &href
		r.Path, _ = url.Parse(href)

		GetRequest(r)
	})
	doc.Find("bgsound").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("body").Each(func(_ int, s *goquery.Selection) {
		background, _ := s.Attr("background")
		//r.Path = &background
		r.Path, _ = url.Parse(background)

		GetRequest(r)
	})
	doc.Find("embed").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		//r.Path = &href
		r.Path, _ = url.Parse(href)

		GetRequest(r)
	})
	doc.Find("embed").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("fig").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		// url.Parseはgetrequest.go, postrequest.goで実装されている
		GetRequest(r)
	})
	doc.Find("frame").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("iframe").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		//r.Path = &href
		r.Path, _ = url.Parse(href)

		GetRequest(r)
	})
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		lowsrc, _ := s.Attr("lowsrc")
		//r.Path = &lowsrc
		r.Path, _ = url.Parse(lowsrc)

		GetRequest(r)
	})
	doc.Find("input").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("input").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("layer").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("object").Each(func(_ int, s *goquery.Selection) {
		data, _ := s.Attr("data")
		//r.Path = &data
		r.Path, _ = url.Parse(data)

		GetRequest(r)
	})
	doc.Find("overlay").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//r.Path = &src
		r.Path, _ = url.Parse(src)

		GetRequest(r)
	})
	doc.Find("table").Each(func(_ int, s *goquery.Selection) {
		background, _ := s.Attr("background")
		//r.Path = &background
		r.Path, _ = url.Parse(background)

		GetRequest(r)
	})
	doc.Find("td").Each(func(_ int, s *goquery.Selection) {
		background, _ := s.Attr("background")
		//r.Path = &background
		r.Path, _ = url.Parse(background)

		GetRequest(r)
	})
	doc.Find("th").Each(func(_ int, s *goquery.Selection) {
		background, _ := s.Attr("background")
		//r.Path = &background
		r.Path, _ = url.Parse(background)

		GetRequest(r)
	})
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
