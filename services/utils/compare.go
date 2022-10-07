package utils

import "golang.org/x/crypto/bcrypt"

// ComparePasswords compares hash and given password
func ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false, err
	}

	return true, nil
}
