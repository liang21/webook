package service

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("1223456#hello")
	generateFromPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	t.Log(string(generateFromPassword))
}
