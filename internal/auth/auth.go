package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func CheckHashPassword(pw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
