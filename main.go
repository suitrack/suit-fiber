package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// create hub for coordination between clients
	hub := NewHub()
	// websocket connection handler
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// create new client with current connection
		NewClient(hub, c)
	}))

	// just listen
	err := app.Listen(":4500")
	if err != nil {
		log.Fatalln(err)
	}
}
