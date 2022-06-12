package repo

import (
	"context"
	"errors"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type DatabaseRepo struct {
	db *mongo.Database
}

var CancellationError = errors.New("cancellation")

const (
	DB_COLLECTION = "users"
)

func NewDatabaseRepo(db *mongo.Database) *DatabaseRepo {
	return &DatabaseRepo{
		db: db,
	}
}

func (r *DatabaseRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	query, err := r.db.Collection(DB_COLLECTION).Find(ctx, bson.D{})
	defer query.Close(ctx)

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)
	for query.Next(ctx) {
		var user models.User
		err := query.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *DatabaseRepo) Get(ctx context.Context, id string) (*models.User, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	query := r.db.Collection(DB_COLLECTION).FindOne(ctx, bson.M{"_id": docId})

	var user models.User
	err = query.Decode(&user)

	return &user, err
}

func (r *DatabaseRepo) GetByName(ctx context.Context, uname string) (*models.User, error) {
	query := r.db.Collection(DB_COLLECTION).FindOne(ctx, bson.M{"username": uname})

	var user models.User
	err := query.Decode(&user)

	return &user, err
}

func (r *DatabaseRepo) Insert(ctx context.Context, user *models.User) error {
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

func (r *DatabaseRepo) Update(ctx context.Context, user *models.User) error {
	dataReq := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"password":   user.Password,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	res := r.db.Collection(DB_COLLECTION).FindOneAndUpdate(ctx, bson.M{"_id": user.ID}, dataReq)

	return res.Err()
}

func (r *DatabaseRepo) UpdatePassword(ctx context.Context, user *models.User) error {
	dataReq := bson.M{
		"$set": bson.M{
			"password": user.Password,
		},
	}

	_, err := r.db.Collection(DB_COLLECTION).UpdateOne(ctx, bson.M{"_id": user.ID}, dataReq)

	return err
}
