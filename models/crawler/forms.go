package crawler

import (
	"net/url"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

var applyData = map[string]string{}

func SetValues(form []entity.HtmlForm, r *entity.RequestStruct) {
	r.Form.Action = form[0].Action
	path, err := url.Parse(form[0].Action)
	if logger.ErrHandle(err) {
		return
	}
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

			case v.Type != "":
				if data, ok := applyData[*v.Name]; ok {
					values.Set(*v.Name, data)
					break
				}
				if v.Placeholder != nil {
					values.Set(*v.Name, *v.Placeholder)
				} else if v.Value != nil {
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
