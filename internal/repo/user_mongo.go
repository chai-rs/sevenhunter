package repo

import (
	"context"
	"net/http"
	"time"

	"github.com/chai-rs/sevenhunter/internal/model"
	errx "github.com/chai-rs/sevenhunter/pkg/error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidUserID = func(e error) error { return errx.E(http.StatusBadRequest, e, "invalid user identification") }
)

type userMongo struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashed_password"`
	CreatedAt      time.Time          `bson:"created_at"`
}

func (u *userMongo) toModel() (*model.User, error) {
	return model.NewUser(model.UserOpts{
		ID:             u.ID.Hex(),
		Name:           u.Name,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		CreatedAt:      u.CreatedAt,
	})
}

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database, collectionName string) *UserRepo {
	return &UserRepo{
		collection: db.Collection(collectionName),
	}
}

var _ model.UserRepo = (*UserRepo)(nil)

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, errx.Mongo(err)
	}

	return count, nil
}

func (r *UserRepo) List(ctx context.Context, opts model.ListUserOpts) ([]model.User, error) {
	filter := bson.M{}
	if opts.Cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(opts.Cursor)
		if err != nil {
			return nil, ErrInvalidUserID(err)
		}

		if opts.SortAsc {
			filter["_id"] = bson.M{"$gt": cursorID}
		} else {
			filter["_id"] = bson.M{"$lt": cursorID}
		}
	}

	findOpts := options.Find()
	findOpts.SetLimit(int64(opts.GetLimit()))

	if opts.SortAsc {
		findOpts.SetSort(bson.D{{Key: "_id", Value: 1}})
	} else {
		findOpts.SetSort(bson.D{{Key: "_id", Value: -1}})
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, errx.Mongo(err)
	}
	defer cursor.Close(ctx)

	var mongoUsers []userMongo
	if err := cursor.All(ctx, &mongoUsers); err != nil {
		return nil, errx.Mongo(err)
	}

	users := make([]model.User, 0, len(mongoUsers))
	for _, mu := range mongoUsers {
		user, err := mu.toModel()
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	u := userMongo{
		Name:           user.Name(),
		Email:          user.Email(),
		HashedPassword: user.HashedPassword(),
		CreatedAt:      user.CreatedAt(),
	}

	result, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return nil, errx.Mongo(err)
	}

	u.ID = result.InsertedID.(primitive.ObjectID)
	return u.toModel()
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidUserID(err)
	}

	var (
		u      userMongo
		filter = bson.M{
			"_id": objID,
		}
	)

	err = r.collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return nil, errx.Mongo(err)
	}

	return u.toModel()
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var (
		u      userMongo
		filter = bson.M{
			"email": email,
		}
	)
	err := r.collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return nil, errx.Mongo(err)
	}

	return u.toModel()
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	objID, err := primitive.ObjectIDFromHex(user.ID())
	if err != nil {
		return ErrInvalidUserID(err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":            user.Name(),
			"email":           user.Email(),
			"hashed_password": user.HashedPassword(),
		},
	}

	filter := bson.M{
		"_id": objID,
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errx.Mongo(err)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidUserID(err)
	}

	filter := bson.M{
		"_id": objID,
	}

	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errx.Mongo(err)
	}

	return nil
}
