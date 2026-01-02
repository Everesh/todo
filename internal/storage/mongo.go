package storage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStorage implements MongoDB storage
type MongoStorage struct {
	collection *mongo.Collection
}

type mongoDocument struct {
	Key  string `bson:"_id"`
	Data []byte `bson:"data"`
}

func NewMongoStorage(uri, database, collection string) (*MongoStorage, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	coll := client.Database(database).Collection(collection)
	return &MongoStorage{collection: coll}, nil
}

func (ms *MongoStorage) Load(key string) ([]byte, error) {
	var doc mongoDocument
	err := ms.collection.FindOne(context.TODO(), bson.M{"_id": key}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("key not found")
		}
		return nil, err
	}
	return doc.Data, nil
}

func (ms *MongoStorage) Save(key string, data []byte) error {
	doc := mongoDocument{
		Key:  key,
		Data: data,
	}
	opts := options.Replace().SetUpsert(true)
	_, err := ms.collection.ReplaceOne(context.TODO(), bson.M{"_id": key}, doc, opts)
	return err
}

func (ms *MongoStorage) Delete(key string) error {
	_, err := ms.collection.DeleteOne(context.TODO(), bson.M{"_id": key})
	return err
}

func (ms *MongoStorage) Exists(key string) (bool, error) {
	count, err := ms.collection.CountDocuments(context.TODO(), bson.M{"_id": key})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
