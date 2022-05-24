package main

import (
	"context"
	"fmt"
	"gitlab.com/g6834/team17/auth_service/internal/db"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func main() {
	mongoCfg := db.MongoConfig{
		Timeout:  5000,
		DBname:   "mts",
		Username: "",
		Password: "",
		Host:     "0.0.0.0",
		Port:     "27017",
	}

	mongo, err := db.Connect(mongoCfg)
	if err != nil {
		log.Fatal(err)
	}

	r := repo.New(mongo)
	id, _ := primitive.ObjectIDFromHex("628bb7b446394be573b1c9dd")
	user := &model.User{
		ID:        id,
		Username:  "PsyMaza",
		Password:  "xxx",
		Email:     "es@ivanov-dev.ru",
		FirstName: "Evgeny",
		LastName:  "Ivanov",
	}

	err = r.Update(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(data)

	fmt.Println("Done")
}
