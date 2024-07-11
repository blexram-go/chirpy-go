package database

import (
	"errors"
	"os"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

var ErrAlreadyExists = errors.New("user already exists")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, pw string) (User, error) {
	if _, err := db.GetUserByEmail(email); !(errors.Is(err, os.ErrNotExist)) {
		return User{}, ErrAlreadyExists
	}
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	newUser := User{
		ID:          id,
		Email:       email,
		Password:    pw,
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = newUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, os.ErrNotExist
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, os.ErrNotExist
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, os.ErrNotExist
	}

	user.Email = email
	user.Password = password
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeChirpyRed(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, os.ErrNotExist
	}

	user.IsChirpyRed = true
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
