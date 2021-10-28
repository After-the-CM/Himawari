package scanner

import (
	"net/http"

	"Himawari/models/entity"
)

func Dirlisting(j *entity.JsonNode) {

	d := determinant{
		kind:          DirList,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	req := createGetReq(j.URL, "")

	//最初は`/`がついていないURL
	d.approach(d, []*http.Request{req})

	reqslash := createGetReq(j.URL+"/", "")

	d.approach(d, []*http.Request{reqslash})

}
