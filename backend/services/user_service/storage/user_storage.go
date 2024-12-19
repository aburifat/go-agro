package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStorage interface {
	Create(user *User) (string, error)
	Update(id string, user *User) error
	GetById(id string) (*User, error)
	GetAll(pageNumber, pageSize int) ([]*User, error)
	Delete(id string) error
}

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type MongoUserStorage struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoUserStorage(uri, dbName string) (*MongoUserStorage, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	database := client.Database(dbName)

	return &MongoUserStorage{client: client, database: database}, nil
}

func (s *MongoUserStorage) Create(user *User) (string, error) {
	collection := s.database.Collection("users")

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (s *MongoUserStorage) Update(id string, user *User) error {
	collection := s.database.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"username": user.Username,
			"email":    user.Email,
		},
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (s *MongoUserStorage) GetById(id string) (*User, error) {
	collection := s.database.Collection("users")
	var user User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}
	return &user, nil
}

func (s *MongoUserStorage) GetAll(pageNumber, pageSize int) ([]*User, error) {
	collection := s.database.Collection("users")
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)

	var users []*User
	cursor, err := collection.Find(context.Background(), bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %v", err)
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return users, nil
}

func (s *MongoUserStorage) Delete(id string) error {
	collection := s.database.Collection("users")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %v", err)
	}
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
