package main

import (
	_ "github.com/dozheiny/it-captal-task/docs/swagger"
	"github.com/dozheiny/it-captal-task/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swag "github.com/gofiber/swagger"
	"log"
)

func main() {
	route := fiber.New()

	// use recover function for recover server when panic comes.
	route.Use(recover.New())
	route.Get("/swagger/*", swag.HandlerDefault) // default

	// Register Routers.
	routes.RegisterAll(route)

	if err := route.Listen(":8080"); err != nil {
		log.Fatalf("Got ERR on listening service: %s", err.Error())
	}
}
