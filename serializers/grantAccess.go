package serializers

import "go.mongodb.org/mongo-driver/bson/primitive"

type GrantAccess struct {
	IDs []primitive.ObjectID `json:"ids"`
}
