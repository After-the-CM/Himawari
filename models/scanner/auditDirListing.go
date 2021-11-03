package scanner

import (
	"fmt"
	"net/http"
	"os"

	"Himawari/models/entity"
)

func auditDirListing(j *entity.JsonNode) {

	d := determinant{
		kind:          dirListing,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	req, err := createGetReq(j.URL, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	//最初は`/`がついていないURL
	d.approach(d, []*http.Request{req})

	reqslash, err := createGetReq(j.URL+"/", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	d.approach(d, []*http.Request{reqslash})

}
