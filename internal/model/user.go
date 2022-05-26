package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"password"`
	PasswordHash string             `bson:"password_hash"`
	Email        string             `json:"email" bson:"email"`
	FirstName    string             `json:"firstName" bson:"first_name"`
	LastName     string             `json:"lastName" bson:"last_name"`
	CreationDate uint64             `json:"creationDate" bson:"creation_date"`
}
