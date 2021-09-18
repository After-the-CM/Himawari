package crawler

import (
	"net/url"

	"Himawari/models/entity"
)

var testData = map[string]string{
	"email":    "Himawari@example.com",
	"url":      "http://example.com",
	"tel":      "00012345678",
	"date":     "2020-12-16",
	"text":     "Himawari",
	"textarea": "Himawari",
	"input":    "I am Himawari",
}

func SetValues(form []entity.HtmlForm, r *entity.RequestStruct) {
	r.Form.Action = form[0].Action
	path, _ := url.Parse(form[0].Action)
	r.Path = path
	r.Form.Method = form[0].Method

	values := url.Values{}
	for _, v := range form {
		if v.Name != nil {
			switch {
			case v.IsOption:
				if len(v.Options) == 0 {
					values.Set(*v.Name, v.Options[0])
				} else {
					values.Set(*v.Name, v.Options[1])
				}
			case v.Type != "submit":
				if v.Placeholder != nil {
					values.Set(*v.Name, *v.Placeholder)
				} else if v.Value == nil {
					values.Set(*v.Name, testData[*v.Name])
				} else {
					values.Set(*v.Name, *v.Value)
				}
			}
		}
		if len(values) != 0 {
			r.Param = values
		}
	}
	JudgeMethod(r)
}
