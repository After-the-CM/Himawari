package scanner

import (
	"bufio"
	"fmt"

	"Himawari/models/entity"
)

func XSS(j *entity.JsonNode) {
	d := determinant{
		kind: reflectedXSS,
		// SetHeaderDocumentRootのために一度approachをセット
		approach:      detectReflectedXSS,
		eachVulnIssue: &j.Issue,
	}

	var payload []string
	p := readfile("models/scanner/payload/" + d.kind + ".txt")
	xssPayload := bufio.NewScanner(p)
	for xssPayload.Scan() {
		payload = append(payload, xssPayload.Text())
	}

	// おそらくreflectのみで、randmarkを送信して検証する必要はなさそう。

	if j.Path == "/" {
		for _, v := range j.Children {
			d.jsonMessage = retrieveJsonMessage(&v)
			if d.jsonMessage != nil {
				for _, v := range payload {
					d.setHeaderDocumentRoot(v)
				}
				break
			}
		}

	}

	for i := 0; i < len(j.Messages); i++ {
		d.jsonMessage = &j.Messages[i]
		d.approach = searchRandmark
		tmpCandidate := make([]entity.JsonMessage, 0)
		d.candidate = &tmpCandidate
		d.gatherCandidates(&entity.JsonNodes)

		fmt.Println(j.Path, *d.candidate)

		if len(*d.candidate) != 0 {
			// stored
			d.kind = storedXSS
			d.approach = detectStoredXSS
		} else {
			// reflect
			d.kind = reflectedXSS
			d.approach = detectReflectedXSS
		}

		for _, v := range payload {
			d.setParam(v)
			if len(j.Messages[i].PostParams) != 0 {
				d.setPostHeader(v)
			} else {
				d.setGetHeader(v)
			}
		}
	}
}
