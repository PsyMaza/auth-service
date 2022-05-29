package repo

import (
	"context"
	"errors"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepo interface {
	Get(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByName(ctx context.Context, uname string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, user *model.User) error
	SetToken(ctx context.Context, userID primitive.ObjectID, td *model.TokenDetails) error
	TokenExist(ctx context.Context, token string, types TokenTypes) error
}

type repo struct {
	db *mongo.Database
}

var CancellationError = errors.New("cancellation")

const (
	DB_COLLECTION = "users"
)

func New(db *mongo.Database) UserRepo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetAll(ctx context.Context) ([]*model.User, error) {
	query, err := r.db.Collection(DB_COLLECTION).Find(ctx, bson.D{})
	defer query.Close(ctx)

	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0)
	for query.Next(ctx) {
		var user model.User
		err := query.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *repo) Get(ctx context.Context, id string) (*model.User, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	query := r.db.Collection(DB_COLLECTION).FindOne(ctx, bson.M{"_id": docId})

	var user model.User
	err = query.Decode(&user)

	return &user, err
}

func (r *repo) GetByName(ctx context.Context, uname string) (*model.User, error) {
	query := r.db.Collection(DB_COLLECTION).FindOne(ctx, bson.M{"username": uname})

	var user model.User
	err := query.Decode(&user)

	return &user, err
}

func (r *repo) Insert(ctx context.Context, user *model.User) error {
	dataReq := bson.M{
		"username":      user.Username,
		"password":      user.Password,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"creation_date": time.Now().Unix(),
	}

	_, err := r.db.Collection(DB_COLLECTION).InsertOne(ctx, dataReq)
	if err != nil {
		return err
	}

	return err
}

func (r *repo) Update(ctx context.Context, user *model.User) error {
	dataReq := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	_, err := r.db.Collection(DB_COLLECTION).UpdateOne(ctx, bson.M{"_id": user.ID}, dataReq)
	return err
}

func (r *repo) UpdatePassword(ctx context.Context, user *model.User) error {
	dataReq := bson.M{
		"$set": bson.M{
			"password": user.Password,
		},
	}

	_, err := r.db.Collection(DB_COLLECTION).UpdateOne(ctx, bson.M{"_id": user.ID}, dataReq)

	return err
}

func (r *repo) SetToken(ctx context.Context, userID primitive.ObjectID, td *model.TokenDetails) error {
	dataReq := bson.M{
		"$set": bson.M{
			"token": bson.M{
				"access_token":  td.AccessToken,
				"refresh_token": td.RefreshToken,
				"at_expires":    td.AtExpires,
				"rt_expires":    td.RtExpires,
			},
		},
	}

	_, err := r.db.Collection(DB_COLLECTION).UpdateOne(ctx, bson.M{"_id": userID}, dataReq)
	return err
}

type TokenTypes uint

const (
	Access TokenTypes = iota
	Refresh
)

func (r *repo) TokenExist(ctx context.Context, token string, types TokenTypes) error {
	var dataReq bson.M
	switch types {
	case Access:
		dataReq = bson.M{
			"token.access_token": token,
		}
	case Refresh:
		dataReq = bson.M{
			"token.refresh_token": token,
		}
	}
	res := r.db.Collection(DB_COLLECTION).FindOne(ctx, dataReq)
	return res.Err()
}
