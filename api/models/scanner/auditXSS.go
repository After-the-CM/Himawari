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
					r.landmark = genLandmark()
					r.setHeaderDocumentRoot(strings.Replace(v, "[landmark]", r.landmark, 1))
				}
				break
			}
		}
	}

	for i := 0; i < len(j.Messages); i++ {
		r.jsonMessage = &j.Messages[i]
		s.jsonMessage = &j.Messages[i]
		s.approach = searchLandmark
		if !QuickScan {
			fmt.Println("GatherCandidateeeeeeeeeee")
			s.gatherCandidates()
		}

		if len(s.jsonMessage.Candidate) != 0 {
			// stored
			s.kind = storedXSS
			s.approach = detectStoredXSS

			for _, v := range payloads {
				s.landmark = genLandmark()
				s.setGetParam(strings.Replace(v, "[landmark]", s.landmark, 1))

				s.landmark = genLandmark()
				s.setPostParam(strings.Replace(v, "[landmark]", s.landmark, 1))

				for _, cookie := range j.Cookies {
					s.landmark = genLandmark()
					s.setCookie(cookie, strings.Replace(v, "[landmark]", s.landmark, 1))
				}

				if len(j.Messages[i].PostParams) != 0 {
					s.landmark = genLandmark()
					s.setPostHeader(strings.Replace(v, "[landmark]", s.landmark, 1))
				} else {
					s.landmark = genLandmark()
					s.setGetHeader(strings.Replace(v, "[landmark]", s.landmark, 1))
				}
			}
		} else {
			// reflect
			r.kind = reflectedXSS
			r.approach = detectReflectedXSS
			for _, v := range payloads {
				r.landmark = genLandmark()
				r.setGetParam(strings.Replace(v, "[landmark]", r.landmark, 1))

				r.landmark = genLandmark()
				r.setPostParam(strings.Replace(v, "[landmark]", r.landmark, 1))

				for _, cookie := range j.Cookies {
					r.landmark = genLandmark()
					r.setCookie(cookie, strings.Replace(v, "[landmark]", r.landmark, 1))
				}

				if len(j.Messages[i].PostParams) != 0 {
					r.landmark = genLandmark()
					r.setPostHeader(strings.Replace(v, "[landmark]", r.landmark, 1))
				} else {
					r.landmark = genLandmark()
					r.setGetHeader(strings.Replace(v, "[landmark]", r.landmark, 1))
				}
			}
		}
	}
}
