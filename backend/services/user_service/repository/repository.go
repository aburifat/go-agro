package repository

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func serialize[T any](source, destination *T) {
	srcValue := reflect.ValueOf(source).Elem()
	dstValue := reflect.ValueOf(destination).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.Field(i)

		if !srcField.IsZero() && dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
}

func Create[T interface{}](c *mongo.Collection, model *T) (string, error) {
	result, err := c.InsertOne(context.Background(), model)
	if err != nil {
		return "", fmt.Errorf("failed to insert: %v", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to parse inserted ID")
	}

	return id.Hex(), nil
}

func Update[T interface{}](c *mongo.Collection, id string, model *T) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %v", err)
	}

	var existing T
	err = c.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&existing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document with ID %s not found", id)
		}
		return fmt.Errorf("failed to fetch document: %v", err)
	}

	serialize(model, &existing)

	update := bson.M{"$set": existing}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %v", err)
	}

	return nil
}

func GetById[T interface{}](c *mongo.Collection, id string) (*T, error) {
	var item T
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	err = c.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}
	return &item, nil
}

func GetAll[T interface{}](c *mongo.Collection, pageNumber, pageSize int) ([]*T, error) {
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)

	var data []*T
	cursor, err := c.Find(context.Background(), bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode user: %v", err)
		}
		data = append(data, &item)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return data, nil
}

func Delete(c mongo.Collection, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %v", err)
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
