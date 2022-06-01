package request

import "regexp"

type Req struct {
	Url       string `json:"url" binding:"required"`
	Uuid      string `json:"uuid"`
	CreatedAt string `json:"createdAt"`
	Tag       Tag    `jsoin:"tag"`
}

type Tag struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
}

func Key(url string) string {
	if len(url) == 0 {
		return ""
	}
	exp := regexp.MustCompile(`^http.*watch\?v\=`)
	if exp.MatchString(url) {
		return exp.ReplaceAllString(url, "")
	}
	return url
}
