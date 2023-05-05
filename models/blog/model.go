package blog

import (
	"github.com/dozheiny/it-captal-task/models/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionName = "blogs"
)

type Model struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
	User    user.Model         `json:"user" bson:"user"`
}

// GetCollectionName returns collection name of blogs' collection in Database.
func (m Model) GetCollectionName() string {
	return collectionName
}

// GetID returns id of document.
func (m Model) GetID() interface{} {
	return m.ID
}
