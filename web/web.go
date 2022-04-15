package web

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func routes(app *fiber.App, c *Controller) {
	app.Static("/", "./static")
	app.Get("/createwallet", c.CreateWallet)
	app.Get("/getbalance", c.GetBalance)
	app.Post("/send", c.Send)
	app.Get("/getblocks", c.GetBlocks)
	app.Get("/getaddress", c.GetAddress)
}

func Init() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	s := NewService()
	c := NewController(s)
	routes(app, c)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	log.Printf("listening on port :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
