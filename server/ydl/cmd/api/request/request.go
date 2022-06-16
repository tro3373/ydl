package request

import (
	"regexp"

	"go.uber.org/zap/zapcore"
)

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

func (r Req) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("url", r.Url)
	enc.AddString("uuid", r.Uuid)
	enc.AddString("createdAt", r.CreatedAt)
	if err := enc.AddObject("tag", r.Tag); err != nil {
		return err
	}
	return nil
}

func (t Tag) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("titile", t.Title)
	enc.AddString("artist", t.Artist)
	enc.AddString("album", t.Album)
	enc.AddString("genre", t.Genre)
	return nil
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
