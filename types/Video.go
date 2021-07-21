package types

import "html/template"

type Video struct {
	Url          string
	RelUrl			 string
	LbryUrl      string
	OdyseeUrl    string
	ClaimId      string
	Channel      Channel
	Title        string
	ThumbnailUrl template.URL
	Description  template.HTML
	License      string
	Views        int64
	Likes        int64
	Dislikes     int64
	Tags         []string
	RelTime			 string
	Date         string
	Duration		 string
}