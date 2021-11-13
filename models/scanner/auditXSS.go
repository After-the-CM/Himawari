package scanner

import (
	"bufio"
	"fmt"
	"strings"

	"Himawari/models/entity"
)

func auditXSS(j *entity.JsonNode) {
	r := determinant{
		kind:          reflectedXSS,
		approach:      detectReflectedXSS,
		eachVulnIssue: &j.Issue,
	}

	s := determinant{
		kind:          storedXSS,
		approach:      detectStoredXSS,
		eachVulnIssue: &j.Issue,
	}

	var payloads []string
	p := readfile("models/scanner/payload/" + "XSS" + ".txt")
	xssPayloads := bufio.NewScanner(p)
	for xssPayloads.Scan() {
		payloads = append(payloads, xssPayloads.Text())
	}

	if j.Path == "/" {
		for _, v := range j.Children {
			r.jsonMessage = retrieveJsonMessage(&v)
			if r.jsonMessage != nil {
				for _, v := range payloads {
					r.setHeaderDocumentRoot(v)
				}
				break
			}
		}
	}

	for i := 0; i < len(j.Messages); i++ {
		r.jsonMessage = &j.Messages[i]
		s.jsonMessage = &j.Messages[i]
		s.approach = searchRandmark
		tmpCandidate := make([]entity.JsonMessage, 0)
		s.candidate = &tmpCandidate
		s.gatherCandidates(&entity.JsonNodes)

		fmt.Println(j.Path, *s.candidate)

		if len(*s.candidate) != 0 {
			// stored
			s.kind = storedXSS
			s.approach = detectStoredXSS

			for _, v := range payloads {
				s.randmark = genRandmark()
				s.setGetParam(strings.Replace(v, "[randmark]", s.randmark, 1))

				s.randmark = genRandmark()
				s.setPostParam(strings.Replace(v, "[randmark]", s.randmark, 1))

				for _, cookie := range j.Cookies {
					s.randmark = genRandmark()
					s.setCookie(cookie, strings.Replace(v, "[randmark]", r.randmark, 1))
				}

				//if fullscan{}
				//scannerutil.gatherCandidates
				/*
					if len(j.Messages[i].PostParams) != 0 {
						s.setPostHeader(v)
					} else {
						s.setGetHeader(v)
					}
				*/
			}
		} else {
			// reflect
			r.kind = reflectedXSS
			r.approach = detectReflectedXSS
			for _, v := range payloads {
				r.randmark = genRandmark()
				r.setGetParam(strings.Replace(v, "[randmark]", r.randmark, 1))

				r.randmark = genRandmark()
				r.setPostParam(strings.Replace(v, "[randmark]", r.randmark, 1))

				for _, cookie := range j.Cookies {
					s.randmark = genRandmark()
					s.setCookie(cookie, strings.Replace(v, "[randmark]", r.randmark, 1))
				}

				if len(j.Messages[i].PostParams) != 0 {
					r.randmark = genRandmark()
					r.setPostHeader(strings.Replace(v, "[randmark]", r.randmark, 1))
				} else {
					r.randmark = genRandmark()
					r.setGetHeader(strings.Replace(v, "[randmark]", r.randmark, 1))
				}
			}
		}
	}
}
