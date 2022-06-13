package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id" yaml:"id"`
	Username     string             `bson:"username" json:"username" yaml:"username"`
	Password     string             `bson:"password" json:"password" yaml:"password"`
	Email        string             `bson:"email" json:"email" yaml:"email"`
	FirstName    string             `bson:"first_name" json:"first_name" yaml:"first_name"`
	LastName     string             `bson:"last_name" json:"last_name" yaml:"last_name"`
	CreationDate uint64             `bson:"creation_date" json:"creationDate"`
}
