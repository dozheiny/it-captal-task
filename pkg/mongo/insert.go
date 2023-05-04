package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert gives model interface with Recorder methods and insert record into DB;
// when a job is done, it returns ID of inserted record.
func (b *Base) Insert(model Recorder) (*primitive.ObjectID, error) {
	b, err := b.SetCollectionName(model.GetCollectionName())
	if err != nil {
		return nil, err
	}

	insertResult, err := b.Collection.InsertOne(b.Context, model)
	if err != nil {
		return nil, err
	}

	id, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, cannotConvertInsertID
	}

	return &id, nil
}

// Insert is a shortcut for base.Insert(model Recorder) .
func Insert(c context.Context, model Recorder) error {
	base, err := Init(c).SetDatabase()
	if err != nil {
		return err
	}

	if _, err := base.Insert(model); err != nil {
		return err
	}

	return nil
}
