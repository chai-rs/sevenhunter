package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongo struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashed_password"`
	CreatedAt      time.Time          `bson:"created_at"`
}

type UserRepo struct {
	collection *mongo.Collection
}
