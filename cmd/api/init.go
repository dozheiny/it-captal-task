package main

import (
	"context"
	"errors"
	"github.com/dozheiny/it-captal-task/database"
	endpoints2 "github.com/dozheiny/it-captal-task/models/endpoints"
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/pkg/redis"
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

	// check Endpoints are existed or not,
	// if not exist, then add it.
	endpoints := make([]endpoints2.Model, 0)

	endpoints = append(endpoints,
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/login", Method: "POST"},
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/refresh", Method: "PATCH"},
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/logout", Method: "PATCH"},
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/blog", Method: "POST"},
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/blog", Method: "GET"},
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/user", Method: "POST"},
		endpoints2.Model{ID: primitive.NewObjectID(), Endpoint: "/user/:user/grant-access", Method: "PUT"},
	)

	for number := range endpoints {
		if err := mongo.FindOne(ctx, endpoints2.Model{}, bson.D{{Key: "endpoint", Value: endpoints[number].Endpoint},
			{Key: "method", Value: endpoints[number].Method}}); err != nil {
			if errors.Is(err, mongo.NoDocuments) {
				if err := mongo.Insert(ctx, endpoints[number]); err != nil {
					log.Fatalf("Gotcha err on save document: %s", err)
				}
			}
		}
	}

	// check admin user exist or not;
	// if not exist, then create a new admin user.
	u := new(user.Model)

	if err := mongo.FindOne(ctx, u, bson.D{{Key: "username", Value: "admin"}}); err != nil {
		if errors.Is(err, mongo.NoDocuments) {
			u.ID = primitive.NewObjectID()
			u.Username = "admin"
			u.SetPassword("admin")
			u.Endpoints = endpoints

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
