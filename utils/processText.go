package utils

import (
	"html"
	"strings"

	"github.com/spf13/viper"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

func ProcessText(text string) string {
	text = bluemonday.UGCPolicy().Sanitize(text)
	text = string(markdown.ToHTML([]byte(text), nil, nil))
	text = strings.ReplaceAll(text, `img src="`, `img src="/image?url=`)
	text = html.UnescapeString(text)

	return text
}

func LbryTo(link string, linkType string) string {
	switch linkType {
	case "rel":
		link = strings.ReplaceAll(link, "lbry://", "/")
		link = strings.ReplaceAll(link, "#", ":")
	case "http":
		link = strings.ReplaceAll(link, "lbry://", "https://" + viper.GetString("DOMAIN") + "/")
		link = strings.ReplaceAll(link, "#", ":")
	case "odysee":
		link = strings.ReplaceAll(link, "lbry://", "https://" + viper.GetString("DOMAIN") + "/")
		link = strings.ReplaceAll(link, "#", ":")
	}
	
	return link
}