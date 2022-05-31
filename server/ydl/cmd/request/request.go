package request

import "regexp"

type Req struct {
	Url       string `json:"url" binding:"required"`
	Tag       Tag    `json:"tag"`
	CreatedAt string `json:"createdAt"`
	Uuid      string `json:"uuid"`
}

type Tag struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
}

func (r Req) Key() string {
	if len(r.Url) == 0 {
		return ""
	}
	exp := regexp.MustCompile(`^http.*watch\?v\=`)
	if exp.MatchString(r.Url) {
		return exp.ReplaceAllString(r.Url, "")
	}
	return r.Url
}
