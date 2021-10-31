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

	var reflectedPayloads []string
	rp := readfile("models/scanner/payload/" + r.kind + ".txt")
	rPayloads := bufio.NewScanner(rp)
	for rPayloads.Scan() {
		reflectedPayloads = append(reflectedPayloads, rPayloads.Text())
	}

	var storedPayloads []string
	sp := readfile("models/scanner/payload/" + s.kind + ".txt")
	sPayloads := bufio.NewScanner(sp)
	for sPayloads.Scan() {
		storedPayloads = append(storedPayloads, sPayloads.Text())
	}

	if j.Path == "/" {
		for _, v := range j.Children {
			r.jsonMessage = retrieveJsonMessage(&v)
			if r.jsonMessage != nil {
				for _, v := range reflectedPayloads {
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

			for _, v := range storedPayloads {
				s.randmark = genRandmark()
				s.setGetParam(strings.Replace(v, "[randmark]", s.randmark, 1))

				s.randmark = genRandmark()
				s.setPostParam(strings.Replace(v, "[randmark]", s.randmark, 1))

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
			for _, v := range reflectedPayloads {
				r.setGetParam(v)
				r.setPostParam(v)

				if len(j.Messages[i].PostParams) != 0 {
					r.setPostHeader(v)
				} else {
					r.setGetHeader(v)
				}
			}
		}
	}
}
