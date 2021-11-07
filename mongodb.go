package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

type MongoLoader interface {
	Load(docs []interface{})
}

func openConnection(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (m *Mongo) InsertMany(database string, collection string, docs []interface{}) {
	coll := m.Client.Database(database).Collection(collection)

	_, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
}

func (m *Mongo) Close() {
	if err := m.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func NewMongo(uri string) *Mongo {
	return &Mongo{
		Client: openConnection(uri),
	}
}
