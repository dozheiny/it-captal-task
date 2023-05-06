package endpoints

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	collectionName = "endpoints"
)

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Endpoint  string             `json:"endpoint" bson:"endpoint"`
	Method    string             `json:"method" bson:"method"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m Model) GetCollectionName() string {
	return collectionName
}

func (m Model) GetID() interface{} {
	return m.ID
}
