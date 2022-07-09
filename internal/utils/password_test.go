package utils_test

import (
	"gitlab.com/g6834/team17/auth-service/internal/utils"
	"testing"
)

func FuzzCheckPassword(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		hash := utils.GetHash(data)
		if err := utils.CheckPassword(data, []byte(hash)); err != nil {
			panic(err)
		}
	})
}
