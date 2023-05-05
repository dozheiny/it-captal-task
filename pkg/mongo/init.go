package mongo

import (
	"context"

	"github.com/dozheiny/it-captal-task/config"
	"github.com/dozheiny/it-captal-task/database"
	"go.mongodb.org/mongo-driver/bson"
)

// Init will make initialize mongoDB client for future usages.
func Init(ctx context.Context) *Base {

	base := new(Base)
	base.Client = database.Client
	base.Context = ctx

	return base
}

// SetDatabase given a DB_NAME and connect to a database
// if there's no database or DB_NAME is wrong, it returns error.
func (b *Base) SetDatabase() (*Base, error) {

	DB_NAME, err := config.Get("DB_NAME")
	if err != nil {
		return nil, err
	}

	names, err := b.Client.ListDatabaseNames(b.Context, bson.D{})
	if err != nil {
		return nil, err
	}

	DB_NAMEs := make(map[string]bool)
	for number := range names {
		DB_NAMEs[names[number]] = true
	}

	if DB_NAMEs[DB_NAME] {
		b.Database = b.Client.Database(DB_NAME)
		b.Client = b.Database.Client()
		b.DataBaseName = DB_NAME
		return b, nil
	}

	return nil, DB_NAMENotFound
}

// SetCollectionName sets collection name for future requests,
// if collection name doesn't exist, it creates collection name.
func (b *Base) SetCollectionName(collectionName string) (*Base, error) {
	names, err := b.Client.Database(b.DataBaseName).ListCollectionNames(b.Context, bson.D{})
	if err != nil {
		return nil, err
	}

	collectionNames := make(map[string]bool)
	for number := range names {
		collectionNames[names[number]] = true
	}

	if collectionNames[collectionName] {
		b.Collection = b.Client.Database(b.DataBaseName).Collection(collectionName)
		b.Database = b.Collection.Database()
		b.Client = b.Database.Client()
		return b, nil
	}

	if err := b.Database.CreateCollection(b.Context, collectionName); err != nil {
		return nil, err
	}

	b.Collection = b.Client.Database(b.DataBaseName).Collection(collectionName)
	b.Database = b.Collection.Database()
	b.Client = b.Database.Client()
	return b, nil
}
