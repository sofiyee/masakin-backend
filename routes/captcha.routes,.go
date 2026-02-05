package routes

import (
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"

	"github.com/fogleman/gg"
	"math/rand"
	"time"
)

func RegisterCaptchaRoutes(app *fiber.App) {

	// 1) Generate MATH captcha (JSON)
	api := app.Group("/api")
	api.Get("/generate-math-captcha", func(c *fiber.Ctx) error {
		id, question := service.GenerateMathCaptcha()

		return c.JSON(fiber.Map{
			"captcha_id": id,
			"question":   question,
			"image_url":  "/captcha/math/" + id + ".png",
		})
	})

	// 2) Render IMAGE captcha (PNG)
	app.Get("/captcha/math/:captchaId.png", func(c *fiber.Ctx) error {
		question := service.GetMathCaptchaQuestion(c.Params("captchaId"))
		if question == "" {
			return c.SendStatus(404)
		}

		// canvas
		width := 200
		height := 100
		dc := gg.NewContext(width, height)

		dc.SetRGB(1, 1, 1)
		dc.Clear()

		_ = dc.LoadFontFace("assets/fonts/DejaVuSans-Bold.ttf", 36)

		text := question + " = ?"

		tw, th := dc.MeasureString(text)

		
		x := float64(width)/2 - tw/2
		y := float64(height)/2 + th/2

		dc.SetRGB(0, 0, 0)
		dc.DrawString(text, x, y)

		
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 6; i++ {
			dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
			dc.SetLineWidth(2)
			dc.DrawLine(
				rand.Float64()*float64(width),
				rand.Float64()*float64(height),
				rand.Float64()*float64(width),
				rand.Float64()*float64(height),
			)
			dc.Stroke()
		}

		c.Set("Content-Type", "image/png")
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		return dc.EncodePNG(c.Response().BodyWriter())
	})
}
