package utils

import (
	"bytes"
	"html"
	"net/url"
	"regexp"
	"strings"

	"codeberg.org/librarian/librarian/data"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func ProcessText(text string, newline bool) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		panic(err)
	}
	text = buf.String()
	if newline {
		text = strings.ReplaceAll(text, "\n\n", "")
		text = strings.ReplaceAll(text, "\n", "<br>")
	}
	re := regexp.MustCompile(`(?:img src=")(.*)(?:")`)
	imgs := re.FindAllString(text, len(text)/4)
	for i := 0; i < len(imgs); i++ {
		hmac := EncodeHMAC(imgs[i])
		text = re.ReplaceAllString(text, "/image?url=$1"+hmac)
	}
	text = strings.ReplaceAll(text, `img src="`, `img src="/image?url=`)
	text = strings.ReplaceAll(text, "https://odysee.com", viper.GetString("DOMAIN"))
	text = strings.ReplaceAll(text, "https://open.lbry.com", viper.GetString("DOMAIN"))
	text = html.UnescapeString(text)
	text = bluemonday.UGCPolicy().RequireNoReferrerOnLinks(true).RequireNoFollowOnLinks(true).RequireCrossOriginAnonymous(true).Sanitize(text)
	text = ReplaceStickersAndEmotes(text)

	return text
}

func ProcessDocument(text string, isMd bool) string {
	if isMd {
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
		)
		var buf bytes.Buffer
		if err := md.Convert([]byte(text), &buf); err != nil {
			panic(err)
		}
		text = buf.String()
	}

	re := regexp.MustCompile(`(?:img src=")(.*)(?:")`)
	imgs := re.FindAllString(text, len(text)/4)
	for i := 0; i < len(imgs); i++ {
		imgUrlRe := regexp.MustCompile(`https?:\/\/[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_:%&;\?\#\/.=]+`)
		imgUrl := imgUrlRe.FindString(imgs[i])
		imgUrl = url.QueryEscape(imgUrl)
		hmac := EncodeHMAC(imgUrl)
		text = strings.ReplaceAll(text, imgs[i], `img src="/image?`+`hash=`+hmac+`&url=`+imgUrl+`"`)
	}

	re2 := regexp.MustCompile(`<iframe src="http(.*)>`)
	text = re2.ReplaceAllString(text, "")

	re3 := regexp.MustCompile(`<iframe src="(.*)>`)
	embeds := re3.FindAllString(text, len(text)/4)
	for i := 0; i < len(embeds); i++ {
		embed := embeds[i]
		newEmbed := strings.ReplaceAll(embed, "#", ":")
		newEmbed = strings.ReplaceAll(newEmbed, "lbry://", "/embed/")
		text = strings.ReplaceAll(text, embed, newEmbed)
	}

	text = strings.ReplaceAll(text, "https://odysee.com", viper.GetString("DOMAIN"))
	text = strings.ReplaceAll(text, "https://open.lbry.com", viper.GetString("DOMAIN"))

	p := bluemonday.UGCPolicy()
	p.AllowImages()
	p.AllowElements("iframe")
	p.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("src").OnElements("iframe")
	p.RequireNoReferrerOnLinks(true)
	p.RequireNoFollowOnLinks(true)
	p.RequireCrossOriginAnonymous(true)
	text = p.Sanitize(text)

	return text
}

func LbryTo(link string) (map[string]string, error) {
	link = strings.ReplaceAll(link, "#", ":")
	split := strings.Split(strings.ReplaceAll(link, "lbry://", ""), "/")
	link = "lbry://" + url.PathEscape(split[0])
	if len(split) > 1 {
		link = "lbry://" + url.PathEscape(split[0]) + "/" + url.PathEscape(split[1])
	}

	link = strings.ReplaceAll(link, "lbry://", "http://domain.tld/")
	parsedLink, err := url.Parse(link)
	if err != nil {
		return map[string]string{}, err
	}
	link = parsedLink.String()

	link = strings.ReplaceAll(link, "%3A", ":")
	link = strings.ReplaceAll(link, "+", "%2B")

	return map[string]string{
		"rel":    strings.ReplaceAll(link, "http://domain.tld/", "/"),
		"http":   strings.ReplaceAll(link, "http://domain.tld/", viper.GetString("DOMAIN")+"/"),
		"odysee": strings.ReplaceAll(link, "http://domain.tld/", "https://odysee.com/"),
	}, nil
}

func UrlEncode(link string) (string, error) {
	link2, err := url.Parse(link)
	return link2.String(), err
}

func ReplaceStickersAndEmotes(text string) string {
	re := regexp.MustCompile(":(.*?):")
	emotes := re.FindAllString(text, len(text)/4)
	for i := 0; i < len(emotes); i++ {
		emote := strings.ReplaceAll(emotes[i], ":", "")
		if data.Stickers[emote] != "" {
			proxiedImage := "/image?width=0&height=200&url=" + data.Stickers[emote] + "&hash=" + EncodeHMAC(data.Stickers[emote])
			htmlEmote := `<img loading="lazy" src="` + proxiedImage + `" height="200px">`

			text = strings.ReplaceAll(text, emotes[i], htmlEmote)
		} else if data.Emotes[emote] != "" {
			proxiedImage := "/image?url=" + data.Emotes[emote] + "&hash=" + EncodeHMAC(data.Emotes[emote])
			htmlEmote := `<img loading="lazy" class="emote" src="` + proxiedImage + `" height="24px">`

			text = strings.ReplaceAll(text, emotes[i], htmlEmote)
		}
	}

	return text
}
