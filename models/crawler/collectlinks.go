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
	//Refererをstringから*url.URLに変更
	//nextStruct.Referer = referer.String()
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

func parseForms(doc *goquery.Document, r entity.RequestStruct) {
	doc.Find("form").Each(func(_ int, s *goquery.Selection) {
		form := entity.HtmlForm{}
		form.Action, _ = s.Attr("action")
		form.Method, _ = s.Attr("method")
		var inputs []entity.HtmlForm

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			//tag := "input"

			nameAttr, _ := s.Attr("name")
			form.Name = nameAttr

			typ, _ := s.Attr("type")
			typ = strings.ToLower(typ)
			form.Type = typ

			_, checked := s.Attr("checked")
			if (typ == "radio" || typ == "checkbox") && !checked {
				//return
			}

			//form.Values.Add("Tag", tag)

			value, ok := s.Attr("value")
			if ok {
				form.Value = value
			}

			placeholder, placeholderB := s.Attr("placeholder")
			if placeholderB {
				form.Placeholder = placeholder
			} else {
				form.Placeholder = "NaN"
			}

			/*
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
			*/
			inputs = append(inputs, form)
		})

		//SetValues(form, r)
		SetValues(inputs, r)

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
		//forms = append(forms, form)

	})

	return
}

/*
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

		// formタグがある限り収集したParameterとその値はformsにて保持できている。
		// fmt.Println(fomrs)
		forms = append(forms, form)

	})

	return forms
}

*/
