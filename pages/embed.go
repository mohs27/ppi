package pages

import (
	"fmt"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func EmbedHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=3600")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; connect-src *; media-src * blob:; block-all-mixed-content")

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if err != nil {
		if strings.ContainsAny(err.Error(), "NOT_FOUND") {
			return c.Status(404).Render("errors/notFound", nil)
		}
		return err
	}

	if utils.Contains(viper.GetStringSlice("blocked_claims"), claimData.Id) {
		return c.Status(451).Render("errors/blocked", fiber.Map{
			"claim": claimData,
		})
	}

	if claimData.StreamType == "video" {
		videoStream, err := api.GetStream(claimData.LbryUrl)
		if err != nil {
			return err
		}

		return c.Render("embed", fiber.Map{
			"stream": videoStream,
			"video":  claimData,
		})
	} else {
		return fmt.Errorf("unsupported stream type: " + claimData.StreamType)
	}
}
