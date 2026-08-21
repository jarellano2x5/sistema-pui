package utils

import "golang.org/x/crypto/bcrypt"

func HashPwd(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func HashChk(pwd, hsh string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hsh), []byte(pwd))
	return err == nil
}
