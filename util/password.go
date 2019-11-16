// password related helper function
package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// default password for test use
const DefaultPassword = "123456"

// hash and salt the password
func Password(pwd string) string {
	hashing, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	return string(hashing)
}

func ComparePassword(hash, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err == nil {
		return true
	} else {
		return false
	}
}
