package scanner

import (
	"fmt"
	"net/http"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

func auditDirListing(j *entity.JsonNode) {
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "DirListingの診断を開始しました")
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
