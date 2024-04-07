package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Dorpper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dorpper

	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, int64, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUserByID(context.Context, string, *types.User) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore {
	if dbname == "" {
		dbname = DBNAME
	}
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}

// Function allow dropping of the collection, ruturning an Error if exists
func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- DROPPING USER COLLECITON")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User

	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User

	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, int64, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		count = 0
	}

	var users []*types.User
	err = cur.All(ctx, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err

	}

	if insertedID, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = insertedID
	} else {
		fmt.Printf("Failed:: res.InsertedID is not a primitive.ObjectID: %v", res.InsertedID)
		return nil, errors.New("Update Failed")
	}
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) UpdateUserByID(ctx context.Context, id string, values *types.User) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	_, err = s.coll.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.D{
			{"$set", values},
		})

	if err != nil {
		return nil, err
	}

	return values, err
}
