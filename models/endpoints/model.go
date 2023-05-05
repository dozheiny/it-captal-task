package endpoints

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Endpoint  string             `json:"endpoint" bson:"endpoint"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
