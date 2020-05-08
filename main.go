// +build linux darwin

package main

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/websocket"
	"log"
)

const TextType = 1

func main() {
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// create hub for coordination between clients
	hub := newHub()
	go hub.run()

	// a dummy get endpoint that says hello
	app.Get("/hello/:name?", func(c *fiber.Ctx) {
		name := c.Params("name")
		if name != "" {
			c.Status(fiber.StatusOK).Send("Hello " + name + " !!")
		} else {
			c.Status(fiber.StatusOK).Send("Hello World !!")
		}
	})

	// websocket connection handler
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// create new client with current connection
		newClient(hub, c)
	}))

	// just listen
	err := app.Listen(4500)
	if err != nil {
		log.Fatalln(err)
	}
}
