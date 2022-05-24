package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	Password     string             `bson:"password.go"`
	Email        string             `bson:"email"`
	FirstName    string             `bson:"first_name"`
	LastName     string             `bson:"last_name"`
	CreationDate uint64             `bson:"creation_date"`
}
