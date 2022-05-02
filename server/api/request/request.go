package request

import "regexp"

type Exec struct {
	Url       string `json:"url" binding:"required"`
	Tag       Tag    `json:"tag"`
	CreatedAt string `json:"createdAt"`
}

type Tag struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
}

func (r Exec) Key() string {
	if len(r.Url) == 0 {
		return ""
	}
	exp := regexp.MustCompile(`^http.*watch\?v\=`)
	if exp.MatchString(r.Url) {
		return exp.ReplaceAllString(r.Url, "")
	}
	return r.Url
}
