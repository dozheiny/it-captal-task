package main

import (
	"context"
	"errors"
	"github.com/dozheiny/it-captal-task/database"
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/pkg/redis"
	"github.com/dozheiny/it-captal-task/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

// init function will run before running the main function.
// that helps to health check your environment is healthy or not.
func init() {

	ctx := context.Background()

	if err := database.Connection(ctx); err != nil {
		log.Fatalf("Gotcha Error on make connection to mongo: %s", err.Error())
	}

	if _, err := redis.GetRedis().GetConnection(); err != nil {
		log.Fatalf("Gotcha Error on get connection from redis: %s", err.Error())
	}

	// check admin user exist or not;
	// if not exist, then create a new admin user.
	u := new(user.Model)

	if err := mongo.FindOne(ctx, u, bson.D{{Key: "username", Value: "admin"}}); err != nil {
		if errors.Is(err, mongo.NoDocuments) {
			u.ID = primitive.NewObjectID()
			u.Username = "admin"
			u.SetPassword("admin")

			if err := mongo.Insert(ctx, u); err != nil {
				log.Fatalf("Gotcha Error on creating admin user: %s", err)
			}
			// if error is not NoDocuments then kill service to running.
			// to check err.
		} else {

			log.Fatalf("Gotcha Error on finding user: %s", err)
		}
	}
}

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
