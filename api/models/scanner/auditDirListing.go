package scanner

import (
	"net/http"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

func auditDirListing(j *entity.JsonNode) {

	d := determinant{
		kind:          dirListing,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	req, err := createGetReq(j.URL, "")
	if logger.ErrHandle(err) {
		return
	}

	//最初は`/`がついていないURL
	d.approach(d, []*http.Request{req})

	reqslash, err := createGetReq(j.URL+"/", "")
	if logger.ErrHandle(err) {
		return
	}
	d.approach(d, []*http.Request{reqslash})

}
