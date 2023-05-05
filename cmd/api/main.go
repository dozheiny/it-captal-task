package main

import (
	"github.com/dozheiny/it-captal-task/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

func main() {
	route := fiber.New()

	// use recover function for recover server when panic comes.
	route.Use(recover.New())

	// Register Routers.
	routes.RegisterAll(route)

	if err := route.Listen(":8080"); err != nil {
		log.Fatalf("Got ERR on listening service: %s", err.Error())
	}
}
