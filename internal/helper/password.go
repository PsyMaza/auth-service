package helper

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func CheckPassword(pwd, pwdHash []byte) error {
	return bcrypt.CompareHashAndPassword(pwdHash, pwd)
}
