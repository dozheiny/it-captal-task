package user

import (
	"crypto/sha256"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collectionName = "users"

type Model struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password []byte             `bson:"password" json:"-"`
}

// GetCollectionName returns collection name of users' collection in Database.
func (m *Model) GetCollectionName() string {
	return collectionName
}

// GetID returns id of one document. (in bson called "_id").
func (m *Model) GetID() interface{} {
	return m.ID
}

// SetPassword gives password as string;
// Then store hash of password into given struct.
func (m *Model) SetPassword(password string) {
	hash := sha256.Sum256([]byte(password))
	m.Password = hash[:]
}

// VerifyPassword gives pas
func (m *Model) VerifyPassword(password string) bool {
	hash := sha256.Sum256([]byte(password))

	if hex.EncodeToString(m.Password) == hex.EncodeToString(hash[:]) {
		return true
	}

	return false
}
