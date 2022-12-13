package main

import (
	"encoding/json"
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iuribrito/go-notification/config"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/register", func(c *fiber.Ctx) error {
		payload := struct {
			Subscription string `json:"subscription"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.JSON(fiber.Map{"error": true, "message": err.Error()})
		}

		fmt.Println(payload.Subscription)
		return c.JSON(fiber.Map{"subscription": payload.Subscription})
	})

	app.Get("/keys", func(c *fiber.Ctx) error {
		publicKey := config.Config("VAPID_PUBLICKEY")
		privateKey := config.Config("VAPID_PRIVATEKEY")

		return c.JSON(fiber.Map{"publicKey": publicKey, "privateKey": privateKey})
	})

	app.Get("/generate_keys", func(c *fiber.Ctx) error {
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			return c.JSON(fiber.Map{"error": true, "message": err.Error()})
		}

		return c.JSON(fiber.Map{"publicKey": publicKey, "privateKey": privateKey})
	})

	app.Post("/send_notification", func(c *fiber.Ctx) error {

		payload := struct {
			Message string `json:"message"`
			Title string `json:"title"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.JSON(fiber.Map{"error": true, "message": err.Error()})
		}

		subscription := ``
		publicKey := config.Config("VAPID_PUBLICKEY")
		privateKey := config.Config("VAPID_PRIVATEKEY")

		s := &webpush.Subscription{}
		json.Unmarshal([]byte(subscription), s)

		var title string
		var message string
		
		if (payload.Title != "") {
			title = payload.Title
		}  else {
			title = "WeeBet"
		}

		if (payload.Message != "") {
			message = payload.Message
		}  else {
			message = "WeeBet"
		}

		notification := `{"notification":{"title":"` + title + `","body":"` + message + `"}}`

		resp, err := webpush.SendNotification([]byte(notification), s, &webpush.Options{
			Subscriber: "email@email.com",
			VAPIDPublicKey: publicKey,
			VAPIDPrivateKey: privateKey,
		})

		if(err != nil) {
			return c.JSON(fiber.Map{"error": true, "message": err.Error()})
		}

		defer resp.Body.Close()
		return c.SendString("Sent!")
	})

	app.Listen(":4000")
}
