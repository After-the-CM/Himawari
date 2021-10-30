package scanner

import (
	"Himawari/models/entity"
)

func Scan(j *entity.JsonNode) {
	Osci(j)
	DirTrav(j)
	SQLi(j)
	Dirlisting(j)
  XSS(j)

	if len(j.Children) > 0 {
		for i := 0; i < len(j.Children); i++ {
			Scan(&j.Children[i])
		}
	}
}
